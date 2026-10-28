// Package game is a stateful logic for yahtzee
package game

import (
	"math/rand"
	"slices"
)

func (g *Game) getYahtzeeStatus() yahtzeeStatus {
	isYahtzee := slices.Contains(g.diceCount[:], 5)
	isSecondYahtzee := isYahtzee && g.card[Yahtzee] > 0
	isJoker := isSecondYahtzee && g.qualifiesAsJoker()
	return yahtzeeStatus{isYahtzee, isSecondYahtzee, isJoker}
}

// CanMove returns if there are any moves left for the game
func (g *Game) CanMove() bool {
	return slices.Contains(g.card, -1)
}

// GetPhase returns the current phase of the game
func (g *Game) GetPhase() Phase {
	return g.phase
}

// Faces returns a sorted slice of the faces of the dice
func (g *Game) Faces() []int {
	faces := make([]int, 0, NumberOfDice)
	for k, v := range g.diceCount {
		for range v {
			faces = append(faces, k+1)
		}
	}
	return faces
}

// returns if the current selection to reroll is valid
func (g *Game) validateReroll(toRoll FaceCounts) bool {
	for i, v := range toRoll {
		if g.diceCount[i] < v {
			return false
		}
	}
	return true
}

// Roll does the  inital roll
// progresses to reroll phase
func (g *Game) Roll() error {
	if g.phase != PhaseRoll {
		return ErrWrongPhase
	}

	var dice FaceCounts
	for range NumberOfDice {
		dice[rand.Intn(FaceCount)]++
	}
	g.diceCount = dice
	g.rerolls = rerolls
	g.phase = PhaseReroll
	return nil
}

// EndRerolling ends the reroll phase
// progresses to the score selection phase
func (g *Game) EndRerolling() error {
	if g.phase != PhaseReroll {
		return ErrWrongPhase
	}
	g.phase = PhaseSelectScore

	return nil
}

// Reroll rerolls dice of given face values providied
// THIS USES DICE FACE COUNTS
// Progresses to the score selection phase if the number of rerolls left is 0
func (g *Game) Reroll(toRoll FaceCounts) error {
	if g.phase != PhaseReroll {
		return ErrWrongPhase
	}
	if g.rerolls == 0 {
		return ErrNoRolls
	}
	if !g.validateReroll(toRoll) {
		return ErrInvalidDiceRoll
	}
	for i, v := range toRoll {
		g.diceCount[i] -= v
		for range v {
			g.diceCount[rand.Intn(FaceCount)] += 1
		}
	}
	g.rerolls--
	if g.rerolls == 0 {
		g.phase = PhaseSelectScore
	}
	return nil
}

// returns the number of points for a of a kind for a given amount
func ofAKindPoints(counts FaceCounts, faces []int, amount int) int {
	for _, c := range counts {
		if c >= amount {
			return totalValueCount(faces)
		}
	}
	return 0
}

func totalValueCount(faces []int) int {
	total := 0
	for _, f := range faces {
		total += f
	}
	return total
}

func fullhouse(counts FaceCounts) int {
	two := false
	three := false
	for _, v := range counts {
		switch v {
		case 2:
			two = true
		case 3:
			three = true
		}
	}

	if two && three {
		return 25
	}
	return 0
}

func straight(counts FaceCounts, length, points int) int {
	longest, run := 0, 0
	for _, v := range counts {
		if v != 0 {
			run++
		} else {
			run = 0
		}
		longest = max(longest, run)
	}
	if longest >= length {
		return points
	}
	return 0
}

func (g *Game) validatePlacePoints(choice Category) bool {
	return slices.Contains(g.options, choice)

}

// returns if the current roll is a joker
// assumes that there is a non-zero yahtzee scored
func (g *Game) qualifiesAsJoker() bool {
	for k, v := range g.diceCount {
		if v != 0 {
			return g.card[k] != UnusedCategoryValue
		}
	}
	panic(ErrNoDiceRolled)
}

// GetPotentialPointsForCategory returns the points that would be scored for a specific slot if selected
func (g *Game) GetPotentialPointsForCategory(slotChoice Category) int {
	status := g.getYahtzeeStatus()

	switch slotChoice {
	case Tok:
		return ofAKindPoints(g.diceCount, g.Faces(), 3)
	case Fok:
		return ofAKindPoints(g.diceCount, g.Faces(), 4)
	case Fh:
		if status.isJoker {
			return 25
		}
		return fullhouse(g.diceCount)
	case Ss:
		if status.isJoker {
			return 30
		}
		return straight(g.diceCount, 4, 30)
	case Ls:
		if status.isJoker {
			return 40
		}
		return straight(g.diceCount, 5, 40)
	case Yahtzee:
		if status.isYahtzee {
			return 50
		}
		return 0
	case Chance:
		return totalValueCount(g.Faces())
	default:
		return g.diceCount[slotChoice] * (int(slotChoice) + 1)
	}
}

// PlacePoints select the slot that you want to score
// progresses to roll phase if there are turns left, otherwise progresses
// to the end phase
func (g *Game) PlacePoints(slotChoice Category) error {
	if g.phase != PhaseSelectScore {
		return ErrWrongPhase
	}
	if !g.validatePlacePoints(slotChoice) {
		return ErrInvalidOption
	}

	idx := slices.Index(g.options, slotChoice)
	g.options = slices.Delete(g.options, idx, idx+1)

	if g.getYahtzeeStatus().isSecondYahtzee {
		g.card[Yahtzee] += 100
	}

	g.card[slotChoice] = g.GetPotentialPointsForCategory(slotChoice)

	if len(g.options) == 0 {
		g.phase = PhaseEnd
	} else {
		g.phase = PhaseRoll
	}
	return nil
}

func (g *Game) TotalPoints() int {
	total := g.upperTotal()
	if g.UpperBonus() {
		total += UpperBonusAmount
	}
	for i := int(Tok); i <= int(Chance); i++ {
		if g.card[i] != UnusedCategoryValue {
			total += g.card[i]
		}
	}
	return total
}

func (g *Game) upperTotal() int {
	total := 0
	// doesn't use a for range loop to keep its starting point dependent on Ones
	for i := int(Ones); i <= int(Sixes); i++ {
		if g.card[i] != UnusedCategoryValue {
			total += g.card[i]
		}
	}
	return total
}

func (g *Game) UpperBonus() bool {
	return g.upperTotal() >= UpperBonusRequirement

}

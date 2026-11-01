package game

import "math/rand/v2"

const (
	UnusedCategoryValue   = -1
	NumberOfDice          = 5
	rerolls               = 2
	FaceCount             = 6
	UpperBonusRequirement = 63
	UpperBonusAmount      = 35
)

type Category int

const (
	Ones Category = iota
	Twos
	Threes
	Fours
	Fives
	Sixes
	Tok
	Fok
	Fh
	Ss
	Ls
	Yahtzee
	Chance

	CategoryCount
)

type Phase int

const (
	PhaseRoll Phase = iota
	PhaseReroll
	PhaseSelectScore
	PhaseEnd
)

type ScoreCard []int
type FaceCounts [FaceCount]int

type yahtzeeStatus struct {
	isYahtzee       bool
	isSecondYahtzee bool
	isJoker         bool
}

type Game struct {
	phase     Phase
	rerolls   int
	diceCount FaceCounts
	card      ScoreCard
	options   []Category
	rng       *rand.Rand
}

// NewGame creates a new yahtzee game
func NewGame(seed rand.Source) *Game {
	card := make([]int, CategoryCount)
	options := make([]Category, CategoryCount)
	for i := range card {
		options[i] = Category(i)
		card[i] = UnusedCategoryValue
	}
	return &Game{
		phase:     PhaseRoll,
		rerolls:   0,
		diceCount: FaceCounts{},
		card:      card,
		options:   options,
		rng:       rand.New(rand.New(seed)),
	}
}

// GetPointsForCategory returns the points that have been placed in a games category
// returns -1 if there hasn't been a point placed yet
func (g *Game) GetPointsForCategory(c Category) int {
	return g.card[c]
}

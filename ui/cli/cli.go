// Package cli
// A command-line interfact for te game yahtzee
package cli

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"yahtzee/game"
)

const (
	y            = "y"
	n            = "n"
	maxDigits    = 4 // derived from the max possible points 1575
	maxKeyLength = 2 // defined incase want to add key combination inputs
)

// NOTE: currently just a wrapper incase the struct needs to know its key
type phase struct {
	name string
	// phase game.Phase
}

type stdinReader struct{ scanner *bufio.Scanner }

type cardRow struct {
	key      string
	name     string
	category game.Category
}

var categories = [...]cardRow{
	{"1", "Ones", game.Ones},
	{"2", "Twos", game.Twos},
	{"3", "Threes", game.Threes},
	{"4", "Fours", game.Fours},
	{"5", "Fives", game.Fives},
	{"6", "Sixes", game.Sixes},
	{"7", "Tok", game.Tok},
	{"8", "Fok", game.Fok},
	{"9", "Fh", game.Fh},
	{"10", "Ss", game.Ss},
	{"11", "Ls", game.Ls},
	{"12", "Yahtzee", game.Yahtzee},
	{"13", "Chance", game.Chance},
}

var categoryByKey = func() map[string]game.Category {
	m := make(map[string]game.Category, len(categories))
	for _, c := range categories {
		m[c.key] = c.category
	}

	return m
}()

var phases = map[game.Phase]phase{
	game.PhaseRoll:        {"Roll"},
	game.PhaseReroll:      {"Reroll"},
	game.PhaseSelectScore: {"ScoreSelect"},
	game.PhaseEnd:         {"Game Over"},
}

// just gets the max length needed for displaying the size of categories
var maxCatagoryLength = func() int {
	maxLen := 0
	for _, c := range categories {
		if len(c.name) > maxLen {
			maxLen = len(c.name)
		}
	}
	return maxLen
}()

// Run is the main entry point for a gui
func Run() {
	in := newStdInput()
	clearScreen()
	players := getPlayers(in)

	//main game loop
	for players[0].game.GetPhase() != game.PhaseEnd {
		for _, p := range players {
			g := p.game

			clearScreen()
			phaseBumper(phaseToString(g))
			roll(g)
			displayDice(g)

			var err error
			for g.GetPhase() == game.PhaseReroll {
				clearScreen()
				phaseBumper(phaseToString(g))
				if err != nil {
					fmt.Printf("[%v]\n", err)
				}
				displayDice(g)
				err = reroll(g, in)
			}

			err = nil
			for g.GetPhase() == game.PhaseSelectScore {
				clearScreen()
				phaseBumper(phaseToString(g))
				if err != nil {
					fmt.Printf("%v\n", err)
				}
				displayDice(g)
				displayCard(g)
				err = SelectCategory(g, in)
			}
		}
	}

	for _, p := range players {
		fmt.Println("This is just to display the card once the game is over")
		displayCard(p.game)
	}

}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func phaseBumper(bumperText string) {
	fmt.Printf(`--------------------
---%s
--------------------
`,
		bumperText)
}

// Player setup
func getPlayers(in *stdinReader) []Player {
	players := make([]Player, 0, 2)
	for {
		fmt.Println("Do you want to add another player y/n")

		var s string
		for {
			s = in.mustReadLine()
			if s != y && s != n {
				fmt.Println("Please input y/n")
				continue
			} else {
				break
			}
		}
		if s == n {
			break
		}

		fmt.Println("Please enter your name")
		name := in.mustReadLine()

		players = append(players, NewPlayer(name))

	}

	return players

}

// Game flow

func roll(g *game.Game) {
	fmt.Println("rolling...")

	err := g.Roll()
	if err != nil {
		log.Fatal(err)
	}
}

// reroll ...
func reroll(g *game.Game, in *stdinReader) error {
	fmt.Println(`which dice do you want to reroll?
(list of numbers or press enter with nothing to end reroll)`)
	input := in.mustReadLine()
	dice, err := parseRollInput(input, game.FaceCount)
	switch err {
	case ErrSkipReroll:
		g.EndRerolling()
	case ErrRerollTooLong:
		return fmt.Errorf("thats too many inputs")
	case ErrRerollAtoi:
		return fmt.Errorf("only input numbers together with no commas or spaces")
	case ErrInputOutOfBounds:
		return fmt.Errorf("only input numbers from 0 to %d", game.FaceCount)
	case nil:
		err = g.Reroll(dice)
		if errors.Is(err, game.ErrInvalidDiceRoll) {
			return fmt.Errorf("please only select dice that are in your current roll")
		} else if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal(err)
	}

	return nil
}

func SelectCategory(g *game.Game, in *stdinReader) error {
	input := in.mustReadLine()
	category, ok := categoryByKey[input]
	if !ok {
		return fmt.Errorf("please input one of the listed keys")
	}

	if g.GetPointsForCategory(category) != game.UnusedCategoryValue {
		return fmt.Errorf("please input the index of a unused category")
	}

	g.PlacePoints(category)
	return nil

}

// Display
func displayDice(g *game.Game) {
	faces := g.Faces()
	var sb strings.Builder

	sb.WriteString("face:  ")
	for k, v := range faces {
		sb.WriteString(strconv.Itoa(v))
		if k < len(faces)-1 {
			sb.WriteString("|")
		}
	}
	// sb.WriteString("\nindex: ")
	//
	// for i := range len(faces) {
	// 	sb.WriteString(strconv.Itoa(i + 1))
	// 	if i < len(faces)-1 {
	// 		sb.WriteString("|")
	// 	}
	// }
	sb.WriteRune('\n')

	fmt.Print(sb.String())
}

func displayPlayers(players []Player) {
	for k, v := range players {
		fmt.Printf("Player %d: %s\n", k, v.Name)
	}

}

func phaseToString(g *game.Game) string {
	switch g.GetPhase() {
	case game.PhaseEnd:
		return "end"
	case game.PhaseReroll:
		return "reroll"
	case game.PhaseRoll:
		return "roll"
	case game.PhaseSelectScore:
		return "score select"
	default:
		return "Shouldnt be possible, please report this to <githubURL>"
	}
}

func displayCard(g *game.Game) {
	for _, v := range categories {
		fmt.Printf("%*s: %-*s| %s\n", maxKeyLength, v.key, maxCatagoryLength, v.name, pointsForCategoryString(g, v.category))
	}
}

func pointsForCategoryString(g *game.Game, c game.Category) string {
	if g.GetPointsForCategory(c) < 0 {
		return "  "
	}
	return fmt.Sprintf("%-*s", maxDigits, strconv.Itoa(g.GetPointsForCategory(c)))
}

// Input / parsing

func parseRollInput(s string, numberOfFaces int) (game.FaceCounts, error) {
	if s == "" {
		return game.FaceCounts{}, ErrSkipReroll
	}
	if len(s) > 5 {
		return game.FaceCounts{}, ErrRerollTooLong
	}

	input, err := strconv.Atoi(s)
	if err != nil {
		return game.FaceCounts{}, ErrRerollAtoi
	}

	digits := game.FaceCounts{}
	for input > 0 {

		d := input % 10
		if d < 0 || d > numberOfFaces {
			return game.FaceCounts{}, ErrInputOutOfBounds
		}
		input /= 10
		digits[d-1]++ //-1 because array index
	}

	return digits, nil

}

func newStdInput() *stdinReader {
	return &stdinReader{bufio.NewScanner(os.Stdin)}
}

func (r *stdinReader) mustReadLine() string {
	if !r.scanner.Scan() {
		err := r.scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("unexpected end of input")
	}
	return strings.TrimSpace(r.scanner.Text())
}

package cli

import (
	"math/rand/v2"
	"yahtzee/game"
)

type Player struct {
	Name string
	game *game.Game
}

func NewPlayer(name string) Player {
	g := game.NewGame(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	return Player{
		Name: name,
		game: g,
	}

}

package cli

import "yahtzee/game"

type Player struct {
	Name string
	game *game.Game
}

func NewPlayer(name string) Player {
	g := game.NewGame()
	return Player{
		Name: name,
		game: g,
	}

}

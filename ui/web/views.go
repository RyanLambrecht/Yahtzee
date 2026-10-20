package web

import "yahtzee/game"

type PlayerView struct {
	Name             string
	PointsByCategory map[game.Category]int
	TotalPoints      int
}

type GameView struct {
	Players []PlayerView
	Round   int
}

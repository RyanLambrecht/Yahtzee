package web

import (
	"log/slog"
	"yahtzee/game"
)

type PlayerView struct {
	Name             string
	PointsByCategory map[game.Category]int
	TotalPoints      int
	Bonus            int
}

type GameView struct {
	Players            []PlayerView
	CurrentPlayerIndex int
	Round              int
	UpperSection       []Category
	LowerSection       []Category
}

type Category struct {
	Name  string
	Index int
}

var UpperCategories = []Category{
	{"Ones", 0}, {"Twos", 1}, {"Threes", 2},
	{"Fours", 3}, {"Fives", 4}, {"Sixes", 5},
}

var LowerCategories = []Category{
	{"TOK", 6}, {"FOK", 7}, {"FH", 8},
	{"SS", 9}, {"LS", 10}, {"Yahtzee", 11}, {"Chance", 12},
}

func (s *server) buildGameView() GameView {
	slog.Info("building game view")
	var gv GameView
	for _, g := range s.Games {
		gv.Players = append(gv.Players, buildPlayerView(g))
	}
	gv.CurrentPlayerIndex = s.CurrentPlayerIndex
	gv.Round = s.Games[s.CurrentPlayerIndex].Round
	gv.UpperSection = UpperCategories
	gv.LowerSection = LowerCategories
	return gv
}

func buildPlayerView(p player) PlayerView {
	slog.Info("building player")

	pts := make(map[game.Category]int)
	for i := range game.CategoryCount {
		pts[i] = p.Game.GetPointsForCategory(i)
	}
	bonus := game.UnusedCategoryValue
	if p.Game.UpperBonus() {
		bonus = game.UpperBonusAmount
	}

	pv := PlayerView{
		Name:             p.Name,
		PointsByCategory: pts,
		TotalPoints:      p.Game.TotalPoints(),
		Bonus:            bonus,
	}

	return pv
}

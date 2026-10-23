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
	Dice             []int
}

type GameView struct {
	Players            []PlayerView
	CurrentPlayerIndex int
	Round              int
	UpperSection       []Category
	LowerSection       []Category
	UnusedSlot         int
}

type Category struct {
	Name  string
	Index int
}

var UpperCategories = []Category{
	{"Ones", int(game.Ones)}, {"Twos", int(game.Twos)}, {"Threes", int(game.Threes)},
	{"Fours", int(game.Fours)}, {"Fives", int(game.Fives)}, {"Sixes", int(game.Sixes)},
}

var LowerCategories = []Category{
	{"T.O.K.", 6}, {"F.O.K.", 7}, {"F.H.", 8},
	{"S.S.", 9}, {"L.S.", 10}, {"Yahtzee", 11}, {"Chance", 12},
}

func (s *server) buildGameView() GameView {
	slog.Info("building game view")

	playerViews := make([]PlayerView, 0, len(s.Games))
	for _, g := range s.Games {
		playerViews = append(playerViews, buildPlayerView(g))
	}

	return GameView{
		Players:            playerViews,
		CurrentPlayerIndex: s.CurrentPlayerIndex, //TDOO: can probably change this to a player view if the gohtml is updated
		Round:              s.Games[s.CurrentPlayerIndex].Round,
		UpperSection:       UpperCategories,
		LowerSection:       LowerCategories,
		UnusedSlot:         game.UnusedCategoryValue,
	}
}

func buildPlayerView(p player) PlayerView {
	slog.Info("building player")

	pts := make(map[game.Category]int)
	for i := range game.CategoryCount {
		pts[i] = p.Game.GetPointsForCategory(i)
	}

	bonus := -game.UpperBonusRequirement
	if !p.Game.UpperBonus() {
		for _, v := range UpperCategories {
			points := p.Game.GetPointsForCategory(game.Category(v.Index))
			if points != game.UnusedCategoryValue {
				bonus += points
			}
		}
	} else {
		bonus = game.UpperBonusAmount
	}

	pv := PlayerView{
		Name:             p.Name,
		PointsByCategory: pts,
		TotalPoints:      p.Game.TotalPoints(),
		Bonus:            bonus,
		Dice:             p.Game.Faces(),
	}

	return pv
}

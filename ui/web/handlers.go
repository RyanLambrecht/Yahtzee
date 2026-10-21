package web

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"yahtzee/game"
)

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	slog.Info("index view", "addr", r.RemoteAddr)
	s.reset()
	view(w, r, "startAddingPlayers", nil)
}

func addPlayersForm(w http.ResponseWriter, r *http.Request) {
	slog.Info("prompting for players", "addr", r.RemoteAddr)
	fragment(w, "addPlayersForm", nil)
}

func nameField(w http.ResponseWriter, r *http.Request) {
	fragment(w, "nameField", nil)
}

func (s *server) turn(w http.ResponseWriter, r *http.Request) {
	slog.Info("starting turn", "player", s.currentPlayer())
	fragment(w, "board", s.buildGameView())

}

// TODO: implement this and add a variant for board that just displays
// the dice, and turn avalible slots for the current card into a form
func (s *server) score(w http.ResponseWriter, r *http.Request) {
	slog.Info("letting player score", "player", s.currentPlayer())

}

func (s *server) reroll(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	var indices []int
	for _, v := range r.PostForm["reroll"] {
		i, err := strconv.Atoi(v)
		if err != nil || i < 0 || i >= game.NumberOfDice {
			http.Error(w, "bad reroll index", http.StatusBadRequest)
			return
		}
		indices = append(indices, i)
	}

	p := s.currentPlayer()
	faces := p.Game.Faces()

	fc := game.FaceCounts{}
	for _, idx := range indices {
		fc[faces[idx]-1]++
	}

	slog.Info("rerolling", "indices", indices, "player", p.Name)
	if err := p.Game.Reroll(fc); err != nil {
		slog.Error("reroll failed", "player", p.Name, "error", err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	fragment(w, "rerollDice", buildPlayerView(p))
}

func (s *server) roll(w http.ResponseWriter, r *http.Request) {
	slog.Info("rolling for player", "playerIndex", s.CurrentPlayerIndex)
	p := s.currentPlayer()
	if err := p.Game.Roll(); err != nil {
		slog.Error("There was a error rolling the players dice")
	}
	fragment(w, "rerollDice", buildPlayerView(p))

}

func (s *server) addPlayers(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	var names []string
	for _, n := range r.PostForm["name"] {
		n = strings.TrimSpace(n)
		if n != "" && len(n) < 30 {
			names = append(names, n)
		}
	}

	if len(names) == 0 {
		http.Error(w, "need at least one player", http.StatusBadRequest)
		return
	}

	s.Mux.Lock()
	for _, n := range names {
		s.Games = append(s.Games, player{
			Name: n,
			Game: game.NewGame(),
		})
	}
	s.Mux.Unlock()
	gameView := s.buildGameView()

	slog.Info("players added", "count", len(names), "names", names)

	view(w, r, "board", gameView)
}

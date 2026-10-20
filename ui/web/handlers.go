package web

import (
	"log/slog"
	"net/http"
	"strings"
	"yahtzee/game"
)

func index(w http.ResponseWriter, r *http.Request) {
	slog.Info("index view", "addr", r.RemoteAddr)
	view(w, r, "startAddingPlayers", nil)
}

func addPlayersForm(w http.ResponseWriter, r *http.Request) {
	slog.Info("prompting for players", "addr", r.RemoteAddr)
	fragment(w, "addPlayersForm", nil)
}

func nameField(w http.ResponseWriter, r *http.Request) {
	fragment(w, "nameField", nil)
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
	defer s.Mux.Unlock()
	for _, n := range names {
		s.Games = append(s.Games, player{
			Name: n,
			Game: game.NewGame(),
		})
	}
	slog.Info("players added", "count", len(names), "names", names)

	view(w, r, "board", nil)
}

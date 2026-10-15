package web

import (
	"net/http"
	"strings"
	"yahtzee/game"
)

func addPlayersForm(w http.ResponseWriter, r *http.Request) {
	render(w, addPlayersTmpl, nil)
}

func nameField(w http.ResponseWriter, r *http.Request) {
	renderNamed(w, addPlayersTmpl, "nameField", nil)
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

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, n := range names {
		s.games = append(s.games, player{
			Name: n,
			Game: game.NewGame(),
		})
	}
	render(w, addPlayersTmpl, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	render(w, indexTmpl, nil)
}

// Package web
// Serves the yahtzee game
package web

import (
	"embed"
	"log"
	"net/http"
	"sync"
	"yahtzee/game"
)

type player struct {
	Name string
	Game *game.Game
}

type server struct {
	Mux                sync.Mutex
	CurrentPlayerIndex int
	Games              []player
}

//go:embed  static
var staticFS embed.FS

func Run() {
	ser := &server{}
	mux := ser.routes()

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(staticFS))
	mux.HandleFunc("GET /{$}", index)
	mux.HandleFunc("GET /addPlayersForm", addPlayersForm)
	mux.HandleFunc("GET /name-field", nameField)
	mux.HandleFunc("POST /addPlayers", s.addPlayers)
	return mux
}

func (s *server) reset() {
	s.Games = []player{}
	s.CurrentPlayerIndex = 0
}

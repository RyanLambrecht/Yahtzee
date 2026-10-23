// Package web
// Serves the yahtzee game
package web

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"yahtzee/game"
)

type player struct {
	Name  string
	Game  *game.Game
	Round int
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

	slog.Info("started up server")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(staticFS))
	mux.HandleFunc("GET /{$}", s.index)
	mux.HandleFunc("GET /addPlayersForm", addPlayersForm)
	mux.HandleFunc("GET /name-field", nameField)
	mux.HandleFunc("POST /addPlayers", s.addPlayers)
	mux.HandleFunc("GET /startRound", s.turn)
	mux.HandleFunc("GET /roll", s.roll)
	mux.HandleFunc("POST /reroll", s.reroll)
	mux.HandleFunc("GET /score", s.score)
	mux.HandleFunc("POST /score/{index}", s.placePointInSlot)
	return mux
}

func (s *server) reset() {
	s.Games = []player{}
	s.CurrentPlayerIndex = 0
}

func (s *server) currentPlayer() player {
	return s.Games[s.CurrentPlayerIndex]
}

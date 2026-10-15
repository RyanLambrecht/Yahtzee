// Package web
// used to run the server for yahtzee
package web

import (
	"embed"
	"html/template"
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
	mu    sync.Mutex
	games []player
}

// embeds the files into the binary
var assets embed.FS

var (
	indexTmpl      = template.Must(template.ParseFS(assets, "templates/index.html"))
	addPlayersTmpl = template.Must(template.ParseFS(assets, "templates/forms/addPlayers.html"))
)

func Run() {
	ser := &server{}
	mux := ser.routes()

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", index)
	mux.HandleFunc("GET /addPlayersForm", addPlayersForm)
	mux.HandleFunc("GET /name-field", nameField)
	mux.HandleFunc("POST /addPlayers", s.addPlayers)
	return mux
}

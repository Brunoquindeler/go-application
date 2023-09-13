package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const (
	leagueRoute = "/league/"
	playerRoute = "/players/"
)
const jsonContentType = "application/json"
const PlayerDB = "player.db"

type playerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() []Player
}

type PlayerServer struct {
	store playerStore
	sync.Mutex
	http.Handler
}

func NewPlayerServer(store playerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle(leagueRoute, http.HandlerFunc(p.leagueHandler))
	router.Handle(playerRoute, http.HandlerFunc(p.playersHandler))

	p.Handler = router

	return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, playerRoute)

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.Lock()
	defer p.Unlock()
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

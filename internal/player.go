package internal

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const playerRoute = "/players/"
const playerDB = "player.db"

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}

type PlayerServer struct {
	store PlayerStore
	sync.Mutex
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	return &PlayerServer{
		store: store,
	}
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

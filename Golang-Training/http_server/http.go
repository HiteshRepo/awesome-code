package http_server

import (
	"fmt"
	"net/http"
	"strings"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type PlayerServer struct {
	store PlayerStore
}

func GetNewPlayerServer(store PlayerStore) *PlayerServer {
	return &PlayerServer{store: store}
}


func (ps *PlayerServer) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	router := http.NewServeMux()
	router.Handle("/players/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		player := strings.TrimPrefix(r.URL.Path, "/players/")

		switch request.Method {
		case http.MethodGet:
			ps.showScore(w, player)
		case http.MethodPost:
			ps.processWin(w, player)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)

		}
	}))
	router.ServeHTTP(w, r)
}

func (ps *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := ps.store.GetPlayerScore(player)
	if len(score) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, _ = fmt.Fprint(w, score)
}

func (ps *PlayerServer) processWin(w http.ResponseWriter, player string) {
	ps.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}


type PlayerStore interface {
	GetPlayerScore(name string) string
	RecordWin(name string)
	GetWins() []string
}

type StubPlayerStore struct {
	scores map[string]string
	winCalls []string
}

func GetNewPlayerStore(scores map[string]string) PlayerStore {
	return &StubPlayerStore{scores: scores}
}

func (s *StubPlayerStore) GetPlayerScore(name string) string {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetWins() []string {
	return s.winCalls
}




package main

import (
	"github.com/hiteshpattanayak-tw/golang-training/http_server"
	"log"
	"net/http"
)

func main() {
	scores := map[string]string{
		"Pepper": "20",
		"Floyd":  "10",
	}

	playerStore := http_server.GetNewPlayerStore(scores)
	playerServer := http_server.GetNewPlayerServer(playerStore)

	handler := http.HandlerFunc(playerServer.ServerHTTP)
	log.Fatal(http.ListenAndServe(":5000", handler))
}

package http_server_test

import (
	"github.com/hiteshpattanayak-tw/golang-training/http_server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_PlayerServer(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/players/Floyd", nil)
	response := httptest.NewRecorder()

	scores := map[string]string{
		"Pepper": "20",
		"Floyd":  "10",
	}

	store := http_server.GetNewPlayerStore(scores)

	playerServer := http_server.GetNewPlayerServer(store)

	playerServer.ServerHTTP(response, request)

	got := response.Body.String()
	want := "10"

	assert.Equal(t, want, got)
}

func Test_PlayerServerNotFound(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/players/SomeOneElse", nil)
	response := httptest.NewRecorder()

	scores := map[string]string{
		"Pepper": "20",
		"Floyd":  "10",
	}

	store := http_server.GetNewPlayerStore(scores)

	playerServer := http_server.GetNewPlayerServer(store)

	playerServer.ServerHTTP(response, request)

	got := response.Code
	want := 404

	assert.Equal(t, want, got)
}

func Test_PlayerServerPost(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost, "/players/Floyd", nil)
	response := httptest.NewRecorder()

	scores := map[string]string{
		"Pepper": "20",
		"Floyd":  "10",
	}

	store := http_server.GetNewPlayerStore(scores)

	playerServer := http_server.GetNewPlayerServer(store)

	playerServer.ServerHTTP(response, request)

	got := response.Code
	want := 202

	assert.Equal(t, want, got)
	assert.Equal(t, len(store.GetWins()), 1)
}

func Test_PlayerServerNotAllowed(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPut, "/players/SomeOneElse", nil)
	response := httptest.NewRecorder()

	scores := map[string]string{
		"Pepper": "20",
		"Floyd":  "10",
	}

	store := http_server.GetNewPlayerStore(scores)

	playerServer := http_server.GetNewPlayerServer(store)

	playerServer.ServerHTTP(response, request)

	got := response.Code
	want := 405

	assert.Equal(t, want, got)
}
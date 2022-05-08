package handlers

import (
	"net/http"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
}

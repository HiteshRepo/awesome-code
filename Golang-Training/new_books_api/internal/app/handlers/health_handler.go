package handlers

import (
	"net/http"
)

type HealthHandler struct {}

func (hh *HealthHandler) HandlerFunc(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

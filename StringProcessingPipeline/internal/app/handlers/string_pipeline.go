package handlers

import (
	"encoding/json"
	"github.com/hiteshrepo/StringProcessingPipeline/internal/app/models"
	"github.com/hiteshrepo/StringProcessingPipeline/internal/pkg/queue"
	"net/http"
)

type StringPipelineHandler struct {
	q           *queue.Queue
}

func GetNewStringPipelineHandler(q *queue.Queue) *StringPipelineHandler {
	return &StringPipelineHandler{q: q}
}

func (sph *StringPipelineHandler) AddStringToQueueHandler(w http.ResponseWriter, r *http.Request) {
	var spr models.StringPipelineRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&spr); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	sph.q.Push(spr.Data)
	respondWithJSON(w, http.StatusAccepted, `{"status": "data processed"}`)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

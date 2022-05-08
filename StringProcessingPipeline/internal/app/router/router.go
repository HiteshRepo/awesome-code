package router

import (
	"github.com/gorilla/mux"
	"github.com/hiteshrepo/StringProcessingPipeline/internal/app/handlers"
)

func GetNewRouter(sph *handlers.StringPipelineHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/sph/", sph.AddStringToQueueHandler).Methods("POST")

	return r
}

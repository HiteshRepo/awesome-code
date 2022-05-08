package router

import (
	"github.com/gorilla/mux"
	"github.com/hiteshrepo/StringProcessingPipeline/internal/app/handlers"
)

func GetNewRouter(sph *handlers.StringPipelineHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/sph", sph.AddStringToQueueHandler).Methods("POST")
	r.HandleFunc("/sph/stop", sph.StopServerHandler).Methods("POST")

	return r
}

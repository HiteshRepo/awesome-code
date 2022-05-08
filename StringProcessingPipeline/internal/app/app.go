package app

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hiteshrepo/StringProcessingPipeline/internal/app/service"
	"github.com/hiteshrepo/StringProcessingPipeline/internal/pkg/queue"
	"log"
	"net/http"
	"time"
)

type App struct {
	q           *queue.Queue
	trimChan    chan string
	capitalChan chan string
	reverseChan chan string
	displayChan chan string

	service     *service.StringPipelineService
	Router   *mux.Router
}

func NewApp(q *queue.Queue, trimChan, capitalChan, reverseChan, displayChan chan string, r *mux.Router) *App {
	srv := service.NewStringPipelineService(q, trimChan, capitalChan, reverseChan, displayChan)
	return &App{q: q, trimChan: trimChan, capitalChan: capitalChan, reverseChan: reverseChan, displayChan: displayChan, service: srv, Router: r}
}

func (a *App) Start(ctx context.Context) {
	go a.service.Start(ctx)

	srv := &http.Server{
		Handler:      a.Router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("starting server")
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
}

func (a *App) Shutdown(cancelFunc context.CancelFunc) {
	fmt.Println("shutting down: closing channels")
	cancelFunc()
	close(a.trimChan)
	close(a.capitalChan)
	close(a.reverseChan)
	close(a.displayChan)
}

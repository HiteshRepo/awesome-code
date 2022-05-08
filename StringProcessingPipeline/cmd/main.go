package main

import (
	"context"
	"github.com/hiteshrepo/StringProcessingPipeline/internal/app"
	"github.com/hiteshrepo/StringProcessingPipeline/internal/app/handlers"
	"github.com/hiteshrepo/StringProcessingPipeline/internal/app/router"
	"github.com/hiteshrepo/StringProcessingPipeline/internal/pkg/queue"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	trimChan := make(chan string)
	capitalChan := make(chan string)
	reverseChan := make(chan string)
	displayChan := make(chan string)

	q := queue.NewQueue()
	sph := handlers.GetNewStringPipelineHandler(q)

	r := router.GetNewRouter(sph)

	ctx, cancel := context.WithCancel(context.Background())
	a := app.NewApp(q, trimChan, capitalChan, reverseChan, displayChan, r)
	a.Start(ctx)

	<-interrupt()

	a.Shutdown(cancel)
}

func interrupt() chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	return interrupt
}
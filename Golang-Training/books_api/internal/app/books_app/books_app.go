package books_app

import (
	"github.com/gorilla/mux"
	"github.com/hiteshpattanayak-tw/golangtraining/books_api/internal/pkg/db/migrations"
	"log"
	"net/http"
	"time"
)

type BooksApp struct {
	Migrator *migrations.Migrator
	Router   *mux.Router
}

func (app *BooksApp) Start() {
	app.Migrator.RunMigrations()

	srv := &http.Server{
		Handler:      app.Router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (app *BooksApp) Shutdown() {

}

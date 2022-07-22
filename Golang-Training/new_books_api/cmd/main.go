package main

import (
	"fmt"
	"github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/app/repositories"
	"github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/app/router"
	"github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/app/services"
	"github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/pkg/configs"
	"github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/pkg/db"
	"github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/pkg/db/migrations"
	"log"
	"net/http"
	"time"
)

func main() {
	appConfig, err := configs.ProvideAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.ProvideDBConn(&appConfig.DBConfig)
	if err != nil {
		log.Fatal(err)
	}

	migrator, err := migrations.ProvideMigrator(appConfig.DBConfig, dbConn)
	if err != nil {
		log.Fatal(err)
	}

	migrator.RunMigrations()

	bookRepo := repositories.GetNewBookRepository(dbConn)
	bookSrv := services.GetNewBooksService(bookRepo)
	r := router.ProvideRouter(bookSrv)

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%d", appConfig.ServerConfig.Host, appConfig.ServerConfig.Port),
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Starting server")
	log.Fatal(srv.ListenAndServe())
}

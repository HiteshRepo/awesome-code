package main

import (
	"github.com/hiteshpattanayak-tw/golangtraining/books_api/internal/app/books_app"
	"github.com/hiteshpattanayak-tw/golangtraining/books_api/internal/app/handlers"
	"github.com/hiteshpattanayak-tw/golangtraining/books_api/internal/app/repositories"
	"github.com/hiteshpattanayak-tw/golangtraining/books_api/internal/app/router"
	"github.com/hiteshpattanayak-tw/golangtraining/books_api/internal/app/service"
	"github.com/hiteshpattanayak-tw/golangtraining/books_api/internal/pkg/db"
	"github.com/hiteshpattanayak-tw/golangtraining/books_api/internal/pkg/db/migrations"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	databaseConfigs := db.DatabaseConfig{
		Dbname:        "",
		Username:      "",
		Password:      "",
		Host:          "",
		Schema:        "",
		Port:          0,
		LogMode:       false,
		SslMode:       "",
		Connection:    db.ConnectionPool{},
		MigrationPath: "",
	}

	db, err := db.ProvideDatabase(databaseConfigs, "books service")
	if err != nil {
		panic(err)
	}

	migrator, err := migrations.ProvideMigrator(databaseConfigs, db)
	if err != nil {
		panic(err)
	}

	booksRepo := repositories.GetNewBookRepository()

	booksSrv := service.GetNewBooksService(booksRepo)

	booksHandler := handlers.GetNewBooksHandler(booksSrv)
	bookHandler := handlers.GetNewBookHandler(booksSrv)

	r := router.GetNewRouter(booksHandler, bookHandler)

	app := books_app.BooksApp{
		Migrator: migrator,
		Router:   r,
	}

	app.Start()

	<-interrupt()

	app.Shutdown()
}

func interrupt() chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	return interrupt
}

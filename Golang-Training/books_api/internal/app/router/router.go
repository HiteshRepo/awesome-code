package router

import (
	"github.com/gorilla/mux"
	"github.com/hiteshpattanayak-tw/golangtraining/books_api/internal/app/handlers"
)

func GetNewRouter(booksHandler *handlers.BooksHandler, bookHandler *handlers.BookHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/books/", booksHandler.BooksHandlerFunc).Methods("GET")
	r.HandleFunc("/book/{isbn:[0-9]+}/", bookHandler.MutateBookHandlerFunc).Methods("DELETE", "PUT")
	r.HandleFunc("/book/", bookHandler.AddBookHandlerFunc).Methods("POST")

	return r
}

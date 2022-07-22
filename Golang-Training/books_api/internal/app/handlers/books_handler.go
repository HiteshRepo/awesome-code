package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/hiteshpattanayak-tw/golangtraining/books_api/internal/app/models"
	"github.com/hiteshpattanayak-tw/golangtraining/books_api/internal/app/service"
	"net/http"
	"strconv"
)

type BooksHandler struct {
	booksService *service.BooksService
}

func GetNewBooksHandler(service *service.BooksService) *BooksHandler {
	return &BooksHandler{booksService: service}
}

func (bh *BooksHandler) BooksHandlerFunc(w http.ResponseWriter, r *http.Request) {
	books := bh.booksService.GetAllBook()
	respondWithJSON(w, http.StatusOK, books)
}

type BookHandler struct {
	booksService *service.BooksService
}

func GetNewBookHandler(service *service.BooksService) *BookHandler {
	return &BookHandler{booksService: service}
}

func (bh *BookHandler) GetBookHandlerFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	isbnStr := vars["isbn"]
	isbn, _ := strconv.Atoi(isbnStr)

	book := bh.booksService.GetBook(isbn)
	if book == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, book)
}

func (bh *BookHandler) AddBookHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var book *models.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	bh.booksService.AddBook(book)
	respondWithJSON(w, http.StatusAccepted, `{"status": "book created"}`)
}

func (bh *BookHandler) MutateBookHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		bh.UpdateBookHandlerFunc(w, r)
	}

	if r.Method == "DELETE" {
		bh.RemoveBookHandlerFunc(w, r)
	}
}

func (bh *BookHandler) UpdateBookHandlerFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	isbnStr := vars["isbn"]
	isbn, _ := strconv.Atoi(isbnStr)

	var book *models.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	book.ISBN = isbn

	bh.booksService.UpdateBook(isbn, book)
	respondWithJSON(w, http.StatusAccepted, `{"status": "book updated"}`)
}

func (bh *BookHandler) RemoveBookHandlerFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	isbnStr := vars["isbn"]
	isbn, _ := strconv.Atoi(isbnStr)

	bh.booksService.RemoveBook(isbn)
	respondWithJSON(w, http.StatusAccepted, `{"status": "book removed"}`)
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

package service

import (
	"github.com/hiteshpattanayak-tw/golangtraining/books_api/internal/app/models"
	"github.com/hiteshpattanayak-tw/golangtraining/books_api/internal/app/repositories"
)

type BooksService struct {
	booksRepo *repositories.BooksRepository
}

func GetNewBooksService(booksRepo *repositories.BooksRepository) *BooksService {
	return &BooksService{booksRepo: booksRepo}
}

func (bs *BooksService) AddBook(book *models.Book) {
	bs.booksRepo.AddBook(book)
}

func (bs *BooksService) GetBook(isbn int) *models.Book {
	return bs.booksRepo.GetBookByISBN(isbn)
}

func (bs *BooksService) GetAllBook() []*models.Book {
	return bs.booksRepo.GetAllBooks()
}

func (bs *BooksService) RemoveBook(isbn int) {
	bs.booksRepo.RemoveBook(isbn)
}

func (bs *BooksService) UpdateBook(isbn int, book *models.Book) {
	bs.booksRepo.UpdateBook(isbn, book)
}

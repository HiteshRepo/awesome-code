package services

import (
	"fmt"
	"github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/app/db/models"
	"github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/app/repositories"
	"github.com/pkg/errors"
)

type BooksService struct {
	booksRepo *repositories.BookRepository
}

func GetNewBooksService(booksRepo *repositories.BookRepository) BooksService {
	return BooksService{booksRepo: booksRepo}
}

func (bs *BooksService) AddBook(book *models.Book) {
	bs.booksRepo.AddBook(book)
}

func (bs *BooksService) GetBook(isbn int) (*models.Book, error) {
	book := bs.booksRepo.GetBook(isbn)
	if book != nil {
		return book, nil
	}
	return nil, errors.New(fmt.Sprintf("book with isbn %d was not found", isbn))
}

func (bs *BooksService) GetAllBooks() ([]*models.Book, error) {
	books, err := bs.booksRepo.GetAllBooks()
	if err != nil {
		return nil, err
	}
	if len(books) == 0 {
		return nil, errors.New("No books present")
	}
	return books, nil
}

func (bs *BooksService) RemoveBook(isbn int) {
	bs.booksRepo.RemoveBook(isbn)
}

func (bs *BooksService) UpdateBook(book *models.Book) {
	bs.booksRepo.UpdateBook(book)
}

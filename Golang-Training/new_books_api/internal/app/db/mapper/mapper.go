package mapper

import (
	dbModels "github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/app/db/models"
	"github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/app/models"
)

type Mapper interface {
	DBBook(*models.Book) *dbModels.Book
	Book(*dbModels.Book) *models.Book
}

func DBBook(book *models.Book) *dbModels.Book {
	dbBook := &dbModels.Book{
		ISBN:   book.ISBN,
		Name:   book.Name,
		Author: book.Author,
	}

	return dbBook
}

func Book(dbBooks *dbModels.Book) *models.Book {
	book := &models.Book{
		ISBN:   dbBooks.ISBN,
		Name:   dbBooks.Name,
		Author: dbBooks.Author,
	}

	return book
}

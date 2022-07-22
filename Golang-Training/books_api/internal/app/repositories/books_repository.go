package repositories

import "github.com/hiteshpattanayak-tw/golangtraining/books_api/internal/app/models"

type BooksRepository struct {
	books []*models.Book
}

func GetNewBookRepository() *BooksRepository {
	return &BooksRepository{}
}

func (br *BooksRepository) AddBook(book *models.Book) {
	br.books = append(br.books, book)
}

func (br *BooksRepository) GetBookByISBN(isbn int) *models.Book {
	for _,b := range br.books {
		if b.GetISBN() == isbn {
			return b
		}
	}
	return nil
}

func (br *BooksRepository) GetAllBooks() []*models.Book {
	return br.books
}

func (br *BooksRepository) RemoveBook(isbn int) {
	idx := -1
	for i,b := range br.books {
		if b.GetISBN() == isbn {
			idx = i
		}
	}

	if idx > -1 {
		br.books[idx] = br.books[len(br.books)-1]
		br.books = br.books[:len(br.books)-1]
	}
}

func (br *BooksRepository) UpdateBook(isbn int, book *models.Book) {
	idx := -1
	for i,b := range br.books {
		if b.GetISBN() == isbn {
			idx = i
		}
	}

	if idx > -1 {
		b := br.books[idx]
		b.Name = book.Name
		b.Author = book.Author
	}
}

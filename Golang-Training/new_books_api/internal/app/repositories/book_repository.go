package repositories

import (
	"github.com/hiteshpattanayak-tw/golangtraining/new_books_api/internal/app/db/models"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func GetNewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (br *BookRepository) AddBook(book *models.Book) {
	books := []models.Book{
		{Name: book.Name, Author: book.Author, ISBN: book.ISBN},
	}
	br.db.Create(&books)
}

func (br *BookRepository) UpdateBook(book *models.Book) {
	isbn := book.ISBN
	br.db.Where("isbn = ?", isbn).Updates(map[string]interface{}{"name": book.Name, "author": book.Author})

}

func (br *BookRepository) GetBook(isbn int) *models.Book {
	var book *models.Book
	br.db.Where("isbn = ?", isbn).First(book)
	return book
}

func (br *BookRepository) GetAllBooks() ([]*models.Book, error) {
	books := make([]*models.Book, 0)
	err := br.db.Preload("BaseAsset").Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (br *BookRepository) RemoveBook(isbn int) {
	br.db.Delete(&models.Book{}, isbn)
	// br.db.Where("isbn = ?", isbn).Delete(&models.Book{}, isbn)
}
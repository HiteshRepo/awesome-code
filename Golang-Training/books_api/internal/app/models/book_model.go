package models

type Book struct {
	ISBN   int    `json:"ISBN"`
	Name   string `json:"Name"`
	Author string `json:"Author"`
}

func GetNewBook(isbn int, name, author string) *Book {
	return &Book{ISBN: isbn, Name: name, Author: author}
}

func (b *Book) GetISBN() int {
	return b.ISBN
}

func (b *Book) GetName() string {
	return b.Name
}

func (b *Book) GetAuthor() string {
	return b.Author
}

func (b *Book) SetAuthor(author string) {
	b.Author = author
}

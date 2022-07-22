package models

type Book struct {
	ISBN   int    `json:"isbn"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

func GetNewBook(isbn int, name, author string) *Book {
	return &Book{
		ISBN:   isbn,
		Name:   name,
		Author: author,
	}
}

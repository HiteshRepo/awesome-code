package models

type Book struct {
	ISBN   int    `json:"isbn"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

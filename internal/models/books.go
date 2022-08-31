package models

import "github.com/jinzhu/gorm"

// Book
type Book struct {
	gorm.Model
	Title     string `json:"title"`
	Author    string `json:"author"`
	Language  string `json:"language"`
	Category  string `json:"category"`
	ISBN10    string `json:"isbn10" binding:"len=10"`
	ISBN13    string `json:"isbn13" binding:"len=13"`
	PageCount int    `json:"page_count"`
	Rating    int    `json:"rating"`
}

type BookModel struct {
	DB *gorm.DB
}

func (b *BookModel) Insert(title string, author string, language string, category string, isbn10 string, isbn13 string, pagecount int) (int, error) {
	book := Book{
		Title:     title,
		Author:    author,
		Language:  language,
		Category:  category,
		ISBN10:    isbn10,
		ISBN13:    isbn13,
		PageCount: pagecount,
	}
	result := b.DB.Create(&book)
	return int(book.ID), result.Error
}

func (b *BookModel) Get(id int) (*Book, error) {
	var book = &Book{}
	b.DB.First(&book, id)
	return book, nil
}

// This will return the 10 most recently created snippets.
func (b *BookModel) Latest(limit int) ([]*Book, error) {
	return nil, nil
}

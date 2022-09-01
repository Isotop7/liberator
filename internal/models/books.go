package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

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

type Result struct {
	Date  time.Time
	Value int
}

func (b *BookModel) SumPageCount() (int, error) {
	//TODO: We need to check for user assigned books and progress
	result := Result{Value: 0}
	b.DB.Table("books").Select("sum(page_count) as value").Scan(&result)
	if result.Value > 0 {
		return result.Value, nil
	} else {
		return result.Value, ErrSumPageCount
	}
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
	result := b.DB.First(&book, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		} else {
			return nil, result.Error
		}
	}

	return book, nil
}

// This will return the 10 most recently created snippets.
func (b *BookModel) Latest(limit int) ([]*Book, error) {
	var books = []*Book{}

	result := b.DB.Limit(limit).Order("created_at desc").Find(&books)

	if result.Error != nil {
		return nil, result.Error
	}

	return books, nil
}

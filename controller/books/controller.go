package books

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Load() (*Controller, error) {
	db, err := gorm.Open(sqlite.Open("liberator.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Book{})
	return &Controller{db}, nil
}

// Controller for books
type Controller struct {
	db *gorm.DB
}

// Book struct
type Book struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Language  string `json:"language"`
	Category  string `json:"category"`
	ISBN10    string `json:"isbn10" binding:"len=10"`
	ISBN13    string `json:"isbn13" binding:"len=13"`
	PageCount int    `json:"pagecount"`
	Rating    int    `json:"rating"`
}

// Index shows a list of s
// GET /s
func (c *Controller) Index() ([]*Book, error) {
	var books = []*Book{}
	err := c.db.Find(&books)
	if err != nil {
		return books, nil
	} else {
		return nil, errors.New("error getting list of books")
	}
}

// New  page
// GET /s/new
func (c *Controller) New() {}

// Create a new
// POST /s
func (c *Controller) Create(name string, age int) (*Book, error) {
	return &Book{}, nil
}

// Show a
// GET /s/:id
func (c *Controller) Show(id int) (*Book, error) {
	return &Book{}, nil
}

// Update a
// PATCH /s/:id
func (c *Controller) Update(id int, name string, age int) error {
	return nil
}

// Delete a
// DELETE /s/:id
func (c *Controller) Delete(id int) error {
	return nil
}

// Edit book page
// GET /book/:id/edit
func (c *Controller) Edit(id int) (*Book, error) {
	return &Book{}, nil
}

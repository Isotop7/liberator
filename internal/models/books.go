package models

import (
	"database/sql"
	"errors"
	"time"
)

// Book
type Book struct {
	ID        uint         `json:"id"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
	Title     string       `json:"title"`
	Author    string       `json:"author"`
	Language  string       `json:"language"`
	Category  string       `json:"category"`
	ISBN10    string       `json:"isbn10"`
	ISBN13    string       `json:"isbn13"`
	PageCount int          `json:"page_count"`
	Review    string       `json:"review"`
}

type BookModel struct {
	DB *sql.DB
}

func (b *BookModel) Insert(title string, author string, language string, category string, isbn10 string, isbn13 string, pagecount int, review string) (int, error) {
	timestamp := time.Now()

	result, err := b.DB.Exec(`
			INSERT INTO books (created_at, updated_at, title, author, language, category, isbn10, isbn13, page_count, review)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, timestamp, timestamp, title, author, language, category, isbn10, isbn13, pagecount, review)

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (b *BookModel) Get(id int) (*Book, error) {
	book := &Book{}

	row := b.DB.QueryRow(`
		SELECT id, created_at, updated_at, deleted_at, title, author, language, category, isbn10, isbn13, page_count, review
		FROM books
		WHERE id = ?
		`, id)
	err := row.Scan(
		&book.ID,
		&book.CreatedAt,
		&book.UpdatedAt,
		&book.DeletedAt,
		&book.Title,
		&book.Author,
		&book.Language,
		&book.Category,
		&book.ISBN10,
		&book.ISBN13,
		&book.PageCount,
		&book.Review,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return book, nil
}

func (b *BookModel) Latest(limit int) ([]*Book, error) {
	var books = []*Book{}

	rows, err := b.DB.Query(`
		SELECT id, created_at, updated_at, title, author, language, category, isbn10, isbn13, page_count, review
		FROM books
		ORDER BY created_at DESC
		LIMIT ?
		`, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		book := &Book{}
		err := rows.Scan(
			&book.ID,
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.Title,
			&book.Author,
			&book.Language,
			&book.Category,
			&book.ISBN10,
			&book.ISBN13,
			&book.PageCount,
			&book.Review,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func (b *BookModel) Search(query string) ([]*Book, error) {
	var books = []*Book{}

	rows, err := b.DB.Query(`
		SELECT id, created_at, updated_at, title, author, language, category, isbn10, isbn13, page_count, review
		FROM books
		WHERE title like ? OR author like ?
		`, "%"+query+"%", "%"+query+"%")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		book := &Book{}
		err := rows.Scan(
			&book.ID,
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.Title,
			&book.Author,
			&book.Language,
			&book.Category,
			&book.ISBN10,
			&book.ISBN13,
			&book.PageCount,
			&book.Review,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

package models

import (
	"database/sql"
	"errors"
	"time"
)

type BooksUserAssignment struct {
	ID        uint
	UserID    int
	BookID    int
	Status    AssignmentStatus
	PagesRead int
	Rating    int
}

type BooksUserAssignmentModel struct {
	DB *sql.DB
}

type Result struct {
	Value int
}

func (bua *BooksUserAssignmentModel) SumPageCount(userid int) int {
	result := &Result{}
	row := bua.DB.QueryRow(`
		SELECT sum(pages_read) as value 
		FROM books_users_assignment
		WHERE user_id = ?`, userid)

	if errors.Is(row.Err(), sql.ErrNoRows) {
		return 0
	}

	err := row.Scan(&result.Value)

	if err == sql.ErrNoRows {
		return 0
	} else {
		return result.Value
	}
}

func (bua *BooksUserAssignmentModel) GetActiveBooks(userid int) ([]*Book, error) {
	var books = []*Book{}

	rows, err := bua.DB.Query(`
		SELECT b.id, b.created_at, b.updated_at, title, author, language, category, isbn10, isbn13, page_count, review
		FROM books b, books_users_assignment bua
		WHERE b.id == bua.books_id 
		AND bua.user_id = ?
		`, userid)

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

func (bua *BooksUserAssignmentModel) IsCurrentlyAssigned(bookid int) bool {
	bookUserAssignment := &BooksUserAssignment{}

	row := bua.DB.QueryRow(`
		SELECT id, user_id, books_id, status, pages_read, rating
		FROM books_users_assignment
		WHERE (status = ? OR status = ?) 
		AND books_id = ?
		LIMIT 1
	`, Assigned.String(), Active.String(), bookid)

	err := row.Scan(
		&bookUserAssignment.ID,
		&bookUserAssignment.UserID,
		&bookUserAssignment.BookID,
		&bookUserAssignment.Status,
		&bookUserAssignment.PagesRead,
		&bookUserAssignment.Rating,
	)

	if err != nil {
		return false
	}

	if bookUserAssignment.UserID != 0 {
		return true
	} else {
		return false
	}
}

func (bua *BooksUserAssignmentModel) IsCurrentlyAssignedToUser(bookid int, userid int) bool {
	bookUserAssignment := &BooksUserAssignment{}

	row := bua.DB.QueryRow(`
		SELECT id, user_id, books_id, status, pages_read, rating
		FROM books_users_assignment
		WHERE (status = ? OR status = ?) 
		AND books_id = ?
		AND user_id = ?
		LIMIT 1
	`, Assigned.String(), Active.String(), bookid, userid)

	err := row.Scan(
		&bookUserAssignment.ID,
		&bookUserAssignment.UserID,
		&bookUserAssignment.BookID,
		&bookUserAssignment.Status,
		&bookUserAssignment.PagesRead,
		&bookUserAssignment.Rating,
	)

	if err != nil {
		return false
	}

	if bookUserAssignment.UserID == userid {
		return true
	} else {
		return false
	}
}

func (bua *BooksUserAssignmentModel) Assign(bookid int, userid int) (int, error) {
	timestamp := time.Now()

	result, err := bua.DB.Exec(`
		INSERT INTO books_users_assignment (created_at, updated_at, user_id, books_id, status, pages_read, rating)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, timestamp, timestamp, userid, bookid, Assigned, 0, 0)

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (bua *BooksUserAssignmentModel) UpdateAssignmentState(bookid int, userid int, status AssignmentStatus) (int, error) {
	timestamp := time.Now()

	result, err := bua.DB.Exec(`
		UPDATE books_users_assignment
		SET status = ?, updated_at = ?
		WHERE user_id = ? AND books_id = ?
	`, status, timestamp, userid, bookid)

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

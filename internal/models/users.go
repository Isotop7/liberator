package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
	Name           string
	Email          string
	HashedPassword []byte
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	created_at := time.Now()

	_, err = u.DB.Exec(`
		INSERT INTO users (created_at, name, email, hashed_password)
		VALUES (?, ?, ?, ?)
	`, created_at, name, email, hashedPassword)

	if err != nil {
		errMessage := err.Error()
		if strings.Compare(errMessage, "UNIQUE constraint failed: users.email") == 0 {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	user := &User{}

	row := u.DB.QueryRow(`
		SELECT id, name, email, hashed_password
		FROM users
		WHERE email = ?
		LIMIT 1
	`, email)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return int(user.ID), nil
}

func (u *UserModel) Exists(id int) (bool, error) {
	user := &User{}

	row := u.DB.QueryRow(`
		SELECT id, name, email, hashed_password
		FROM users
		WHERE id = ?
		LIMIT 1
	`, id)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrInvalidCredentials
		} else {
			return false, err
		}
	}

	if user.ID != 0 {
		return true, nil
	} else {
		return false, nil
	}
}

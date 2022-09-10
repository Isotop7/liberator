package models

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name           string
	Email          string `gorm:"unique"`
	HashedPassword []byte `gorm:"size:60"`
}

type UserModel struct {
	DB *gorm.DB
}

func (u *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user := User{
		Name:           name,
		Email:          email,
		HashedPassword: hashedPassword,
	}

	result := u.DB.Create(&user)
	if result.Error != nil {
		errMessage := result.Error.Error()
		if strings.Compare(errMessage, "UNIQUE constraint failed: users.email") == 0 {
			return ErrDuplicateEmail
		}
		return result.Error
	}
	return nil
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	user := User{}
	result := u.DB.Where(&User{Email: email}).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, result.Error
		}
	}

	err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return int(user.Model.ID), nil
}

func (u *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

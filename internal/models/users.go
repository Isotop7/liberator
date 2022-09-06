package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name           string
	Email          string `gorm:"unique"`
	HashedPassword []byte `gorm:"size:60"`
}

type UserModel struct {
	DB *gorm.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

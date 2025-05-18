package repository

import (
	"gorm.io/gorm"
)

type User struct {
	// gorm.Model
	// FirstName string `gorm:"not null"`
}

type UserRepository struct {
	DB *gorm.DB
}

func (u UserRepository) Insert() string {
	return "users"
}

func (u UserRepository) Update() string {
	return "users"
}

func (u UserRepository) Delete() string {
	return "users"
}

func (u UserRepository) GetByEmail(email string) interface{} {
	return UserRepository{}
}

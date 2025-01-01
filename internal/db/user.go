package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"uniqueIndex;not null"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Password string `json:"password,-" gorm:"not null"`
}

func NewUser(un, em, pw string) *User {
	return &User{
		Username: un,
		Email:    em,
		Password: pw,
	}
}

func (u *User) Insert() string {
	return "users"
}

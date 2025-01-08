package repository

import (
	"gorm.io/gorm"
)

// AuthModel is a struct that defines the model in the Database
type AuthModel struct {
	gorm.Model
	FirstName string `json:"first_name,omitempty" gorm:"not null"`
	LastName  string `json:"last_name,omitempty" gorm:"not null"`
	Email     string `json:"email,omitempty" gorm:"unique;not null"`
	Password  string `json:"password,omitempty" gorm:"not null"`
	Website   string `json:"website,omitempty" gorm:"not null"`
}

// AuthRepository is a struct that defines the repository for the auth
type AuthRepository struct {
	db *gorm.DB
}

// Create is a method that creates a new auth
func (a AuthRepository) Create(user *AuthModel) error {
	return a.db.Create(&user).Error
}

func (a AuthRepository) SelectByEmail(user *AuthModel) error {
	return a.db.Where("email = ?", user.Email).First(user).Error
}

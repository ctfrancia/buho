package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// ErrRecordNotFound is an error that is returned when a record is not found
var ErrRecordNotFound = errors.New("record not found")

// AuthModel is a struct that defines the model in the Database
type Auth struct {
	gorm.Model
	FirstName          string    `json:"first_name,omitempty" gorm:"not null"`
	LastName           string    `json:"last_name,omitempty" gorm:"not null"`
	Email              string    `json:"email,omitempty" gorm:"unique;not null"`
	Password           string    `json:"password,omitempty" gorm:"not null"`
	Website            string    `json:"website,omitempty" gorm:"not null"`
	RefreshToken       string    `json:"refresh_token,omitempty"`
	RefreshTokenExpiry time.Time `json:"refresh_token_expiry,omitempty"`
}

// AuthRepository is a struct that defines the repository for the auth
type AuthRepository struct {
	db *gorm.DB
}

// Create is a method that creates a new auth
func (a AuthRepository) Create(user *Auth) error {
	return a.db.Create(&user).Error
}

// SelectByEmail is a method that selects a user by email
func (a AuthRepository) SelectByEmail(user *Auth) error {
	return a.db.Where("email = ?", user.Email).First(user).Error
}

// UpdateRefreshToken is a method that updates the refresh token
func (a AuthRepository) Update(user *Auth) error {
	return a.db.Updates(user).Error
}

func (a AuthRepository) SelectByRefreshToken(refreshToken string) (Auth, error) {
	auth := Auth{}
	err := a.db.Where("refresh_token = ?", refreshToken).First(&auth).Error
	if err == gorm.ErrRecordNotFound {
		return auth, ErrRecordNotFound
	}

	return auth, err
}

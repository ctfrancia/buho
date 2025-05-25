package domain

import (
	"time"
)

// LoginRequest represents the request body for the login endpoint
type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

// Auth represents the core domain entity for authentication
type Consumer struct {
	ID                 int64     `json:"id"`
	UUID               string    `json:"uuid"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Email              string    `json:"email"`
	Password           string    `json:"password"`
	Website            string    `json:"website"`
	RefreshToken       string    `json:"refresh_token"`
	RefreshTokenExpiry time.Time `json:"refresh_token_expiry"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

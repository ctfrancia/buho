package entities

import (
	"time"
)

type Player struct {
	ID                 uint
	Name               string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Website            string
	APIKey             string
	Secret             string
	IsActive           bool
	RateLimit          int
	Scopes             string
	RefreshToken       string
	RefreshTokenExpiry time.Time
}

package models

import (
	"github.com/google/uuid"
	"time"
)

// APIClient represents a client application that can access the API
type APIClient struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name      string    `gorm:"not null" json:"name"`          // "Chess Club Mobile App"
	APIKey    string    `gorm:"uniqueIndex;not null" json:"-"` // Don't expose in JSON
	Secret    string    `gorm:"not null" json:"-"`             // Hashed secret
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	RateLimit int       `gorm:"default:1000" json:"rate_limit"` // requests per hour
	Scopes    string    `json:"scopes"`                         // "read:tournaments,write:matches"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

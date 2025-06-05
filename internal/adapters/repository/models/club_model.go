package models

import (
	"github.com/google/uuid"
	"time"
)

// Club represents a chess club
type Club struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Many-to-many relationship
	Players []Player `gorm:"many2many:club_memberships;" json:"players,omitempty"`

	// Junction table access
	ClubMemberships []ClubMembership `json:"club_memberships,omitempty"`
}

package models

import (
	"github.com/google/uuid"
	"time"
)

// Tournament represents a chess tournament
type Tournament struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name            string    `gorm:"not null" json:"name"`
	StartDate       time.Time `gorm:"type:date" json:"start_date"`
	EndDate         time.Time `gorm:"type:date" json:"end_date"`
	Location        string    `json:"location"`
	MaxParticipants int       `json:"max_participants"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// One-to-many relationship
	Matches []Match `gorm:"foreignKey:TournamentID" json:"matches,omitempty"`

	// Many-to-many relationship
	Players []Player `gorm:"many2many:tournament_registrations;" json:"players,omitempty"`

	// Junction table access
	TournamentRegistrations []TournamentRegistration `json:"tournament_registrations,omitempty"`
}

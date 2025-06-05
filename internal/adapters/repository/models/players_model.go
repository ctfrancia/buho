package models

import (
	"github.com/google/uuid"
	"time"
)

// Player represents a chess player
type Player struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Many-to-many relationships
	Clubs       []Club       `gorm:"many2many:club_memberships;" json:"clubs,omitempty"`
	Tournaments []Tournament `gorm:"many2many:tournament_registrations;" json:"tournaments,omitempty"`
	Matches     []Match      `gorm:"many2many:match_participants;" json:"matches,omitempty"`

	// Junction table access
	ClubMemberships         []ClubMembership         `json:"club_memberships,omitempty"`
	TournamentRegistrations []TournamentRegistration `json:"tournament_registrations,omitempty"`
	MatchParticipations     []MatchParticipant       `json:"match_participations,omitempty"`
}

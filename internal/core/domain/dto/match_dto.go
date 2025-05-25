package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateNewMatchRequest struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TournamentID uuid.UUID `gorm:"type:uuid" json:"tournament_id"`
	RoundNumber  int       `json:"round_number"`
	BoardNumber  int       `json:"board_number"`
	MatchDate    time.Time `json:"match_date"`
	Result       string    `json:"result"` // 'white_wins', 'black_wins', 'draw', 'ongoing'
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Foreign key relationship
	// Tournament Tournament `gorm:"foreignKey:TournamentID" json:"tournament,omitempty"`

	// Many-to-many relationship
	//Players []Player `gorm:"many2many:match_participants;" json:"players,omitempty"`

	// Junction table access
	// MatchParticipants []MatchParticipant `json:"match_participants,omitempty"`
}

package entities

import (
	"time"
)

type Match struct {
	ID           uint `gorm:"primarykey"`
	TournamentID uint `gorm:"not null"`
	RoundNumber  int  `gorm:"not null"`
	BoardNumber  int  `gorm:"not null"`
	MatchDate    time.Time
	Result       string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`

	// Foreign key relationship
	Tournament Tournament `gorm:"foreignKey:TournamentID" json:"tournament,omitempty"`

	// Many-to-many relationship
	Players []Player `gorm:"many2many:match_participants;" json:"players,omitempty"`

	// Junction table access
	MatchParticipants []MatchParticipant `json:"match_participants,omitempty"`
}

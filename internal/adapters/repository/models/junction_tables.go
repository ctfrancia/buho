package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// ClubMembership represents the relationship between players and clubs
type ClubMembership struct {
	PlayerID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"player_id"`
	ClubID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"club_id"`
	JoinedDate       time.Time `gorm:"type:date;default:CURRENT_DATE" json:"joined_date"`
	MembershipStatus string    `gorm:"default:'active'" json:"membership_status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Foreign key relationships
	Player Player `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	Club   Club   `gorm:"foreignKey:ClubID" json:"club,omitempty"`
}

// TournamentRegistration represents the relationship between players and tournaments
type TournamentRegistration struct {
	PlayerID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"player_id"`
	TournamentID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"tournament_id"`
	RegistrationDate time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"registration_date"`
	Status           string    `gorm:"default:'registered'" json:"status"` // 'registered', 'confirmed', 'withdrawn'
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Foreign key relationships
	Player     Player     `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	Tournament Tournament `gorm:"foreignKey:TournamentID" json:"tournament,omitempty"`
}

// MatchParticipant represents the relationship between players and matches
type MatchParticipant struct {
	MatchID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"match_id"`
	PlayerID  uuid.UUID `gorm:"type:uuid;primaryKey" json:"player_id"`
	Color     string    `json:"color"` // 'white', 'black'
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Foreign key relationships
	Match  Match  `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	Player Player `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
}

// BeforeCreate hooks for generating UUIDs (if not using database default)
func (p *Player) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (c *Club) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (t *Tournament) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (m *Match) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// TableName methods to explicitly set table names (optional)
func (ClubMembership) TableName() string {
	return "club_memberships"
}

func (TournamentRegistration) TableName() string {
	return "tournament_registrations"
}

func (MatchParticipant) TableName() string {
	return "match_participants"
}

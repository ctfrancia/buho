/*
package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
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

// Match represents a chess match
type Match struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TournamentID uuid.UUID `gorm:"type:uuid" json:"tournament_id"`
	RoundNumber  int       `json:"round_number"`
	BoardNumber  int       `json:"board_number"`
	MatchDate    time.Time `json:"match_date"`
	Result       string    `json:"result"` // 'white_wins', 'black_wins', 'draw', 'ongoing'
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Foreign key relationship
	Tournament Tournament `gorm:"foreignKey:TournamentID" json:"tournament,omitempty"`

	// Many-to-many relationship
	Players []Player `gorm:"many2many:match_participants;" json:"players,omitempty"`

	// Junction table access
	MatchParticipants []MatchParticipant `json:"match_participants,omitempty"`
}

// Junction tables for many-to-many relationships with additional fields

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
*/

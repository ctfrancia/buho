package repository

import (
	"gorm.io/gorm"
	"time"
)

type Tournament struct {
	gorm.Model     `json:"-"`
	Name           string    `gorm:"not null"`
	StartDate      time.Time `gorm:"not null"`
	EndDate        time.Time `gorm:"not null"`
	TournamentUUID string    `gorm:"not null"`
	CreatorID      uint      `json:"-" gorm:"not null"`
	// LocationID  LocationModel
	PosterURL   string
	Description string
	//	Organizer   OrganizerModel
	// Website string
	// Online  bool
}

type TournamentRepository struct {
	DB *gorm.DB
}

func (r TournamentRepository) Create(tournament interface{}) error {
	return r.DB.Create(tournament).Error
}

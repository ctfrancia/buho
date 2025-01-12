package repository

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Tournament struct {
	gorm.Model     `json:"-"`
	Name           string    `json:"name" gorm:"not null"`
	StartDate      time.Time `json:"start_date" gorm:"not null"`
	EndDate        time.Time `json:"end_date" gorm:"not null"`
	TournamentUUID string    `json:"tourament_uuid" gorm:"not null"`
	CreatorID      uint      `json:"-" gorm:"not null"`
	// LocationID  LocationModel
	PosterURL   string `json:"poster_url"`
	Description string `json:"description"`
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

func (r TournamentRepository) GetByUUID(uuid string) (Tournament, error) {
	var tournament Tournament
	err := r.DB.Where("tournament_uuid = ?", uuid).First(&tournament).Error
	return tournament, err
}

func (r TournamentRepository) UpdateByUUID(uuid string, t Tournament) error {
	fmt.Println("Updating tournament", t)
	return r.DB.Model(&Tournament{}).Omit("creator_id", "start_date").Where("tournament_uuid = ?", uuid).Updates(t).Error
}

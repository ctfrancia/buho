package repository

import (
	"time"

	"gorm.io/gorm"
)

type Tournament struct {
	gorm.Model `json:"-"`
	// Name of the tournament
	Name string `json:"name" gorm:"not null"`
	// StartDate of the tournament in the format of time.Time RFC3339
	StartDate time.Time `json:"start_date" gorm:"not null"`
	// EndDate of the tournament in the format of time.Time RFC3339
	EndDate time.Time `json:"end_date" gorm:"not null"`
	// TournamentUUID is the unique identifier of the tournament
	TournamentUUID string `json:"tourament_uuid" gorm:"not null"`
	// CreatorID is the unique identifier of the creator of the tournament
	CreatorID uint `json:"-" gorm:"not null"`
	// PosterURL is the url to the poster. This can be an external link or
	// the path in the sftp server
	PosterURL string `json:"poster_url"`
	// Description of the tournament can be anything the creator wants to add
	Description string `json:"description"`
	// QRCode is the external url or the path in the sftp server
	QRCode string `json:"qr_code"`
}

type TournamentRepository struct {
	DB *gorm.DB
}

func (r TournamentRepository) Create(tournament any) error {
	return r.DB.Create(tournament).Error
}

func (r TournamentRepository) GetByUUID(uuid string) (Tournament, error) {
	var tournament Tournament
	err := r.DB.Where("tournament_uuid = ?", uuid).First(&tournament).Error
	return tournament, err
}

func (r TournamentRepository) UpdateByUUID(uuid string, t Tournament) error {
	return r.DB.Model(&Tournament{}).Omit("creator_id", "start_date").Where("tournament_uuid = ?", uuid).Updates(t).Error
}

func (r TournamentRepository) RemoveTournamentPosterURL(uuid string) error {
	return r.DB.Model(&Tournament{}).Where("tournament_uuid = ?", uuid).Update("poster_url", "").Error
}

func (r TournamentRepository) GetAll() ([]Tournament, error) {
	var tournaments []Tournament
	err := r.DB.Find(&tournaments).Error
	return tournaments, err
}

func (r TournamentRepository) GetByDateRange(startDate, endDate string) ([]Tournament, error) {
	var tournaments []Tournament
	err := r.DB.Where("start_date >= ? AND end_date <= ?", startDate, endDate).Find(&tournaments).Error
	return tournaments, err
}

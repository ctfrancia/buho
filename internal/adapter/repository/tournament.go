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

func (pr *PostgresRepository) CreateTournament(tournament any) error {
	return pr.db.Create(tournament).Error
}

package repository

import (
	"gorm.io/gorm"
)

type TournamentRepository struct {
	DB *gorm.DB
}

func (r TournamentRepository) Create(tournament interface{}) error {
	return r.DB.Create(tournament).Error
}

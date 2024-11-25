package repository

import (
	"gorm.io/gorm"
)

type TournamentsRepository struct {
	DB *gorm.DB
}

func (r TournamentsRepository) Create(tournament interface{}) error {
	return r.DB.Create(tournament).Error
}

package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	Tournaments TournamentsRepository
}

func New(db *gorm.DB) Repository {
	return Repository{
		Tournaments: TournamentsRepository{db},
	}
}

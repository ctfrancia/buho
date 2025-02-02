package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	Tournaments TournamentRepository
	Users       UserRepository
	Auth        AuthRepository
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		Tournaments: TournamentRepository{db},
		Users:       UserRepository{db},
		Auth:        AuthRepository{db},
	}
}

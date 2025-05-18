package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	Tournaments TournamentRepository
	Users       UserRepository
	Auth        AuthRepository
}

// New creates a new repository instance
func New(db *gorm.DB) (Repository, error) {
	return Repository{
		Tournaments: TournamentRepository{db},
		Users:       UserRepository{db},
		Auth:        AuthRepository{db},
	}, nil
}

package repository

import (
	"context"
	"errors"

	domain "github.com/ctfrancia/buho/internal/core/domain/dto"
	"gorm.io/gorm"
)

// GormTournamentRepository is a struct that defines the repository for the tournament
type GormTournamentRepository struct {
	db *gorm.DB
}

// NewGormTournamentRepository creates a new GORM tournament repository
func NewGormTournamentRepository(db *gorm.DB) *GormTournamentRepository {
	return &GormTournamentRepository{
		db: db,
	}
}

// toDomain converts TournamentModel to domain.Tournament
func (r *GormTournamentRepository) toDomain(model TournamentModel) domain.Tournament {
	return domain.Tournament{
		ID:        model.ID,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// toModel converts domain.Tournament to TournamentModel
func (r *GormTournamentRepository) toModel(tournament domain.Tournament) TournamentModel {
	return TournamentModel{
		ID:        tournament.ID,
		Name:      tournament.Name,
		CreatedAt: tournament.CreatedAt,
		UpdatedAt: tournament.UpdatedAt,
	}
}

func (r *GormTournamentRepository) CreateNewTournament(ctx context.Context, tournament domain.Tournament) (domain.Tournament, error) {
	model := r.toModel(tournament)

	result := r.db.WithContext(ctx).Create(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return domain.Tournament{}, secondary.ErrTournamentAlreadyExists
		}
		return domain.Tournament{}, result.Error
	}

	return r.toDomain(model), nil
}

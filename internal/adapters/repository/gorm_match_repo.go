package repository

import (
	"context"
	"github.com/ctfrancia/buho/internal/core/domain/entities"
	"gorm.io/gorm"
)

type GormMatchRepository struct {
	db *gorm.DB
}

func NewGormMatchRepository(db *gorm.DB) *GormMatchRepository {
	return &GormMatchRepository{
		db: db,
	}
}

func (r *GormMatchRepository) SaveMatch(ctx context.Context, match entities.Match) (entities.Match, error) {
	return match, r.db.Create(&match).Error
}

package repository

import (
	"context"
	"errors"
	"time"

	domain "github.com/ctfrancia/buho/internal/core/domain/dto"
	"github.com/ctfrancia/buho/internal/core/domain/entities"
	"gorm.io/gorm"
)

type ConsumerModel struct {
	gorm.Model
	UUID               string    `gorm:"uniqueIndex;not null"`
	FirstName          string    `gorm:"not null"`
	LastName           string    `gorm:"not null"`
	Email              string    `gorm:"uniqueIndex;not null"`
	Password           string    `gorm:"not null"`
	Website            string    `gorm:"not null"`
	RefreshToken       string    `gorm:""`
	RefreshTokenExpiry time.Time `gorm:""`
}

// GormAuthRepository is a struct that defines the repository for the auth
type GormAuthRepository struct {
	db *gorm.DB
}

// NewGormAuthRepository creates a new GORM auth repository
func NewGormAuthRepository(db *gorm.DB) *GormAuthRepository {
	return &GormAuthRepository{
		db: db,
	}
}

// TableName overrides the table name used by GORM
func (ConsumerModel) TableName() string {
	return "consumer"
}

// toDomain converts AuthModel to domain.Auth
func (r *GormAuthRepository) toDomain(model ConsumerModel) domain.Consumer {
	return domain.Consumer{
		UUID:               model.UUID,
		FirstName:          model.FirstName,
		LastName:           model.LastName,
		Email:              model.Email,
		Password:           model.Password,
		Website:            model.Website,
		RefreshToken:       model.RefreshToken,
		RefreshTokenExpiry: model.RefreshTokenExpiry,
	}
}

// toModel converts domain.Auth to AuthModel
func (r *GormAuthRepository) toModel(auth domain.Consumer) ConsumerModel {
	return ConsumerModel{
		UUID:               auth.UUID,
		FirstName:          auth.FirstName,
		LastName:           auth.LastName,
		Email:              auth.Email,
		Password:           auth.Password,
		Website:            auth.Website,
		RefreshToken:       auth.RefreshToken,
		RefreshTokenExpiry: auth.RefreshTokenExpiry,
	}
}

// Create persists a new auth record
func (r *GormAuthRepository) CreateNewApiClient(ctx context.Context, consumer entities.ApiClient) (domain.Consumer, error) {
	model := r.toModel(auth)

	result := r.db.WithContext(ctx).Create(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return domain.Auth{}, secondary.ErrEmailAlreadyExists
		}
		return domain.Auth{}, result.Error
	}

	return r.toDomain(model), nil
}

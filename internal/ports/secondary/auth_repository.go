package secondary

import (
	"context"
	"errors"

	"github.com/ctfrancia/buho/internal/core/domain"
)

// Common repository errors
var (
	ErrAuthNotFound       = errors.New("auth record not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type AuthRepositoryPort interface {
	Create(ctx context.Context, auth domain.Auth) (domain.Auth, error)
	// GetByEmail(ctx context.Context, email string) (domain.Auth, error)
	// Update(ctx context.Context, auth domain.Auth) (domain.Auth, error)
}

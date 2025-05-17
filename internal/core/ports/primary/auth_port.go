package primary

import (
	"context"

	"github.com/ctfrancia/buho/internal/core/domain"
	"github.com/ctfrancia/buho/internal/model"
)

type AuthPort interface {
	// CreateAuthToken(ctx context.Context, req model.CreateAuthTokenRequest) (string, error)
	// NewAPIConsumer(ctx context.Context, req model.NewAPIConsumerRequest) (domain.User, error)
	// RefreshToken(ctx context.Context, req model.RefreshTokenRequest) (string, error)
}

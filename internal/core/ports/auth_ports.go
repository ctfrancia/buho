package ports

import (
	"context"
	"net/http"

	"github.com/ctfrancia/buho/internal/core/domain"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	NewApiConsumer(w http.ResponseWriter, r *http.Request)
}

type AuthService interface {
	// VerifyJWTWithED25519(token string, publicKeyPath string) (*domain.Consumer, error)
}

type AuthRepository interface {
	Create(ctx context.Context, consumer domain.Consumer) (any, error)
}

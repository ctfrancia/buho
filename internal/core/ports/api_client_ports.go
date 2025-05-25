package ports

import (
	"context"
	"net/http"

	"github.com/ctfrancia/buho/internal/core/domain"
)

type ApiClientHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	NewApiConsumer(w http.ResponseWriter, r *http.Request)
}

type ApiClientService interface {
	// VerifyJWTWithED25519(token string, publicKeyPath string) (*domain.Consumer, error)
}

type ApiClientRepository interface {
	CreateNewApiClient(ctx context.Context, consumer domain.Consumer) (any, error)
}

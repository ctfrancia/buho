package ports

import (
	"context"
	"net/http"

	"github.com/ctfrancia/buho/internal/core/domain"
	"github.com/ctfrancia/buho/internal/core/domain/entities"
)

type MatchHandler interface {
	HandleNewMatchRequest(w http.ResponseWriter, r *http.Request)
}

type MatchService interface {
	ProcessNewMatchRequest(ctx context.Context, match entities.Match) (entities.Match, error)
}

type MatchRepository interface {
	SaveMatch(ctx context.Context, match entities.Match) (entities.Match, error)
}

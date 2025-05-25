package ports

import (
	"context"
	"net/http"

	"github.com/ctfrancia/buho/internal/core/domain"
	"github.com/ctfrancia/buho/internal/core/domain/entities"
)

type TournamentHandler interface {
	ProcessNewTournamentRequest(w http.ResponseWriter, r *http.Request)
}

type TournamentService interface {
	CreateNewTournament(ctx context.Context, tournament entities.NewTournamentRequest) (entities.NewTournamentResponse, error)
}

type TournamentRepository interface {
	CreateNewTournament(ctx context.Context, tournament Tournament) (entities.NewTournamentRequest, error)
}

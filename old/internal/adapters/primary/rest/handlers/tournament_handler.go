package handlers

import (
	"github.com/ctfrancia/buho/internal/core/ports"
	"net/http"
)

type TournamentHandler struct {
	tService ports.TournamentService
}

func NewTournamentHandler(tService ports.TournamentService) *TournamentHandler {
	return &TournamentHandler{
		tService: tService,
	}
}

func (h *TournamentHandler) CreateTournament(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

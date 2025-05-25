package handlers

import (
	"github.com/ctfrancia/buho/internal/core/ports"
	"net/http"
)

type TournamentHandler struct {
	service ports.TournamentService
}

func NewTournamentHandler(service ports.TournamentService) *TournamentHandler {
	return &TournamentHandler{
		service: service,
	}
}

func (h *TournamentHandler) ProcessNewTournamentRequest(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
}

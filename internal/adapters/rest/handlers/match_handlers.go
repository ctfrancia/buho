package handlers

import (
	"github.com/ctfrancia/buho/internal/core/ports"
	"net/http"
)

type MatchHandler struct {
	service ports.MatchService
}

func NewMatchHandler(service ports.MatchService) *MatchHandler {
	return &MatchHandler{
		service: service,
	}
}

func (h *MatchHandler) HandleNewMatchRequest(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
}

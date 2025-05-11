package handlers

import (
	"github.com/ctfrancia/buho/internal/core/ports"
	"net/http"
)

type HTTPHandler struct {
	tService ports.TournamentService
	// uService service.UserService
	// eService service.EmailService
}

func NewHTTPHandler(tService ports.TournamentService) *HTTPHandler {
	return &HTTPHandler{
		tService: tService,
	}
}

func (h *HTTPHandler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

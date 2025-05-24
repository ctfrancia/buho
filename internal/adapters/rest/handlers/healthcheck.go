package handlers

import (
	"github.com/ctfrancia/buho/internal/core/ports"
	"net/http"
)

type HealthCheckkHandler struct {
	hch ports.HealthCheckService
}

func NewHealthCheckHandler(hch ports.HealthCheckService) *HealthCheckkHandler {
	return &HealthCheckkHandler{
		hch: hch,
	}
}

func (h *HealthCheckkHandler) Handle(w http.ResponseWriter, r *http.Request) {
	h.hch.GetInformation()
	w.WriteHeader(http.StatusOK)
}

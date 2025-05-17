package handlers

import (
	"net/http"
)

type MiscHandler struct{}

func NewMiscHandler() *MiscHandler {
	return &MiscHandler{}
}

func (h *MiscHandler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

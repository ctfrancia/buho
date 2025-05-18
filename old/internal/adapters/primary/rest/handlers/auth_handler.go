package handlers

import (
	"github.com/ctfrancia/buho/internal/core/ports"
)

type AuthHandler struct {
	aService ports.AuthService
}

func NewAuthHandler(aService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		aService: aService,
	}
}

package repository

import (
	"os"
	"time"

	"github.com/ctfrancia/buho/internal/core/domain/models"
)

// NewServerConfig creates a new server configuration
func NewServerConfig(env string) (*models.Config, error) {
	// TODO: this will take in the details from the cli and structure it here
	return &models.Config{
		Server: models.Server{
			Env:          env,
			Addr:         ":4000",
			IdleTimeout:  2 * time.Minute,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		DB: models.DB{
			DSN: os.Getenv("BUHO_DB_DSN"),
		},
	}, nil
}

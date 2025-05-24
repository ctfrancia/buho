package services

import (
	"github.com/ctfrancia/buho/internal/core/ports"
)

type HealthCheckService struct {
	adapter ports.HealthCheckAdapter
}

func NewHealthCheckService() *HealthCheckService {
	return &HealthCheckService{}
}

func (s *HealthCheckService) GetInformation() {
	// ports.HealthCheckHandler
}

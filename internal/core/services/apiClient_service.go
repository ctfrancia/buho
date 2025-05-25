package services

import (
	"context"
	"github.com/ctfrancia/buho/internal/core/ports"
)

type ConsumerService struct {
	repo ports.AuthRepository
}

func NewApiClientService(repo ports.ApiClientRepository) *ConsumerService {
	return &ConsumerService{
		repo: repo,
	}
}

func (s *ConsumerService) CreateNewAPIConsumer(ctx context.Context) (any, error) {
	return s.repo.Create(ctx, domain.Consumer{})
}

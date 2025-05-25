package services

import (
	"context"
	"github.com/ctfrancia/buho/internal/core/domain"
	"github.com/ctfrancia/buho/internal/core/ports"
)

type ConsumerService struct {
	repo ports.AuthRepository
}

func NewConsumerService(repo secondary.ConsumerRepositoryPort) *ConsumerService {
	return &ConsumerService{
		repo: repo,
	}
}

func (s *ConsumerService) CreateNewAPIConsumer(ctx context.Context) (any, error) {
	return s.repo.Create(ctx, domain.Consumer{})
}

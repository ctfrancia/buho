package services

import (
	"context"
	"github.com/ctfrancia/buho/internal/ports/secondary"
)

type ConsumerService struct {
	repo secondary.ConsumerRepositoryPort
}

func NewConsumerService(repo secondary.ConsumerRepositoryPort) *ConsumerService {
	return &ConsumerService{
		repo: repo,
	}
}

func (s *ConsumerService) CreateNewAPIConsumer(ctx context.Context) (any, error) {
	return s.repo.Create(ctx, domain.Consumer{})
}

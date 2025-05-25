package services

import (
	"context"

	"github.com/ctfrancia/buho/internal/core/domain"
	"github.com/ctfrancia/buho/internal/core/domain/entities"
	"github.com/ctfrancia/buho/internal/core/ports"
)

type MatchService struct {
	repository ports.MatchRepository
}

func NewMatchService(repository ports.MatchRepository) *MatchService {
	return &MatchService{
		repository: repository,
	}
}

func (s *MatchService) ProcessNewMatchRequest(ctx context.Context, match entities.Match) (entities.Match, error) {
	return s.repository.SaveMatch(ctx, match)
}

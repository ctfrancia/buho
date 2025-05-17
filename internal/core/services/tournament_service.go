package services

import (
	"github.com/ctfrancia/buho/internal/core/ports"
)

type TournamentService struct {
	repo ports.TournamentRepository
}

func NewTournamentService(tp ports.TournamentRepository) *TournamentService {
	return &TournamentService{
		repo: tp,
	}
}

func (ts *TournamentService) CreateTournament(tournament any) error {
	return nil
}

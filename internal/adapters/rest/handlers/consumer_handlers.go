package handlers

import (
	"context"

	"github.com/ctfrancia/buho/internal/ports/primary"
)

type ConsumerHandlers struct {
	cService primary.ConsumerServicePort
}

func NewConsumerHandlers(cService primary.ConsumerServicePort) *ConsumerHandlers {
	return &ConsumerHandlers{
		cService: cService,
	}
}

func (h *ConsumerHandlers) CreateNewAPIConsumer(ctx context.Context) (any, error) {
	return h.cService.CreateNewAPIConsumer(ctx)
}

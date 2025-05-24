package handlers

import (
	"context"

	"github.com/ctfrancia/buho/internal/ports/primary"
)

type ConsumerHandlers struct {
	cService primary.ConsumerServicePort
}

func NewConsumerHandler(cService primary.ConsumerServicePort) *ConsumerHandlers {
	return &ConsumerHandlers{
		cService: cService,
	}
}

func (h *ConsumerHandlers) CreateNewAPIConsumer(ctx context.Context) (any, error) {
	return nil, nil
}

func (h *ConsumerHandlers) GetAPIConsumer(ctx context.Context, id string) (any, error) {
	return nil, nil
}

package primary

import (
	"context"
	// "github.com/ctfrancia/buho/internal/core/domain"
)

type ConsumerServicePort interface {
	CreateNewAPIConsumer(ctx context.Context) (any, error)
}

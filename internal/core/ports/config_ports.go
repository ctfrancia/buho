package ports

import "github.com/ctfrancia/buho/old/internal/model"

type ConfigService interface {
	GetConfig() model.Config
}

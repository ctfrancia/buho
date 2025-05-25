package cli

import (
	"flag"

	"github.com/ctfrancia/buho/internal/core/domain/models"
)

func FetchFlags() models.Config {
	isDev := flag.Bool("dev", false, "Run in development mode")

	flag.Parse()

	if *isDev {
		// TODO: Load dev config
		return models.Config{}
	}

	// TODO: Load config
	return models.Config{}
}

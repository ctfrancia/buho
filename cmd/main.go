package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	env string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := config{
		env: "development",
	}

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := http.Server{
		Addr:         ":4000",
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

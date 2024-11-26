package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ctfrancia/buho/internal/model"
	"github.com/ctfrancia/buho/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const version = "1.0.0"

type config struct {
	env string
	db  struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

type application struct {
	config     config
	logger     *slog.Logger
	repository repository.Repository
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var cfg config
	cfg.env = "development"
	cfg.db.dsn = os.Getenv("BUHO_DB_DSN")

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		config:     cfg,
		logger:     logger,
		repository: repository.New(db),
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

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func openDB(cfg config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.db.dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Tournament{})

	return db, nil
}

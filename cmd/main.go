package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ctfrancia/buho/internal/auth"
	"github.com/ctfrancia/buho/internal/digitalocean"
	"github.com/ctfrancia/buho/internal/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	version = "1.0.0"
	// sshAddr is the address of the SSH server local only atm
	sshAddr = "localhost:2022"
)

type application struct {
	config       *Config
	logger       *slog.Logger
	repository   repository.Repository
	auth         *auth.Auth
	digitalOcean *digitalocean.DigitalOceanSpacesClient
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := newConfig()
	cfg.env = "development"
	cfg.db.dsn = os.Getenv("BUHO_DB_DSN")

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	app := &application{
		config:     cfg,
		logger:     logger,
		repository: repository.New(db),
		auth:       auth.NewAuth(cfg.auth.privateKeyPath, cfg.auth.publicKeyPath),
		digitalOcean: digitalocean.NewDigitalOceanSpacesClient(
			cfg.digitalOcean.endpoint,
			cfg.digitalOcean.accessKeyID,
			cfg.digitalOcean.secretAccessKey,
			cfg.digitalOcean.bucket,
		),
	}

	srv := http.Server{
		Addr:         cfg.addr,
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

func openDB(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.db.dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}

	// db.Debug() for debugging
	db.AutoMigrate(&repository.Auth{}, &repository.Organizer{}, &repository.Location{}, &repository.Tournament{})

	return db, nil
}

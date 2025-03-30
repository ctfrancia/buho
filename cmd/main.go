package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ctfrancia/buho/internal/auth"
	"github.com/ctfrancia/buho/internal/digitalocean"
	"github.com/ctfrancia/buho/internal/repository"
	"github.com/ctfrancia/buho/pkg/logger"
	"go.uber.org/zap"

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
	logger       *zap.Logger // *slog.Logger
	repository   repository.Repository
	auth         *auth.Auth
	digitalOcean *digitalocean.DigitalOceanSpacesClient
}

func main() {
	env := os.Getenv("ENV")
	logger := logger.New(env)
	cfg, err := newConfig(env)
	if err != nil {
		logger.Fatal("error creating config", zap.Error(err))
		os.Exit(1)
	}
	cfg.db.dsn = os.Getenv("BUHO_DB_DSN")

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal("error opening db", zap.Error(err))
		os.Exit(1)
	}

	do, err := digitalocean.NewDigitalOceanSpacesClient(
		cfg.digitalOcean.endpoint,
		cfg.digitalOcean.accessKeyID,
		cfg.digitalOcean.secretAccessKey,
		cfg.digitalOcean.bucket,
	)
	if err != nil {
		logger.Fatal("error creating digital ocean client", zap.Error(err))
		os.Exit(1)
	}

	repo, err := repository.New(db)
	if err != nil {
		logger.Fatal("error creating repository", zap.Error(err))
		os.Exit(1)
	}
	auth, err := auth.NewAuth(cfg.auth.privateKeyPath, cfg.auth.publicKeyPath)
	if err != nil {
		logger.Fatal("error creating auth", zap.Error(err))
		os.Exit(1)
	}
	// Create the application
	app := &application{
		config:       cfg,
		logger:       logger,
		repository:   repo,
		auth:         auth,
		digitalOcean: do,
	}

	// Clear out any existing logs when the main function is ending
	defer app.logger.Sync()

	srv := http.Server{
		Addr:         cfg.addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	fields := []zap.Field{
		zap.String("addr", srv.Addr),
		zap.String("env", cfg.env),
	}
	app.logger.Info("starting server", fields...)

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

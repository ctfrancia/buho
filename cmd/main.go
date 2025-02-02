package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ctfrancia/buho/internal/auth"
	// "github.com/ctfrancia/buho/internal/model"
	"github.com/ctfrancia/buho/internal/repository"
	"github.com/ctfrancia/buho/internal/sftp"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	version = "1.0.0"
	// sshAddr is the address of the SSH server local only atm
	sshAddr = "localhost:2022"
)

type application struct {
	config     *config
	logger     *slog.Logger
	repository *repository.Repository
	sftp       *sftp.SSHServer
	auth       *auth.Auth
}

func main() {
	// TODO: Have a config set up function
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// var cfg config
	cfg := NewConfig()

	db, err := openDB(cfg.db.dsn)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	app := &application{
		config:     cfg,
		logger:     logger,
		repository: repository.New(db),
		sftp:       sftp.NewSSHServer("localhost", 2022, "id_rsa", "internal/sftp/pub_key"),
		auth:       auth.NewAuth(cfg.auth.secretKey),
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

func openDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}

	db.Debug().AutoMigrate(&repository.Auth{}, &repository.Organizer{}, &repository.Location{}, &repository.Tournament{})

	return db, nil
}

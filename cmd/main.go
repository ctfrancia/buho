package main

import (
	"fmt"

	"github.com/ctfrancia/buho/internal/model"
	"github.com/ctfrancia/buho/internal/repository"
	// "github.com/pkg/sftp"

	// "github.com/pkg/sftp"
	"github.com/ctfrancia/buho/internal/sftp"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	version = "1.0.0"
	// sshAddr is the address of the SSH server local only atm
	sshAddr = "localhost:2022"
)

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

	sftp := sftp.NewSSHServer("localhost", 2022, "id_rsa")
	sftp.UploadFile()
	// sftp.UploadFile()

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

func HostKeyCb(registeredKey ssh.PublicKey) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		if string(key.Marshal()) == string(registeredKey.Marshal()) {
			return nil
		}

		return fmt.Errorf("host key mismatch")
	}
}

func LoadRegisteredPublicKey(path string) (ssh.PublicKey, error) {
	pubKeyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return pubKey, nil
}

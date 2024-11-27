package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ctfrancia/buho/internal/model"
	"github.com/ctfrancia/buho/internal/repository"
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

	key, err := os.ReadFile("id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key: ", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}
	auth := ssh.PublicKeys(signer)

	// host key callback func
	hostKeyCb := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		// lookup and verify host key...
		fmt.Println("hostname: ", hostname)
		fmt.Println("remote: ", remote)
		fmt.Println("key: ", key)
		return nil
	}

	// ssh client config
	config := &ssh.ClientConfig{
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: hostKeyCb,
	}

	// connect to ssh server
	conn, err := ssh.Dial("tcp", sshAddr, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer conn.Close()

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

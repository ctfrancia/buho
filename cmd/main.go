package main

import (
	"fmt"
	"github.com/ctfrancia/buho/internal/model"
	"github.com/ctfrancia/buho/internal/repository"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"
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

	// buho ssh server
	key, err := os.ReadFile("id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key: ", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}
	auth := ssh.PublicKeys(signer)

	// buho-sftp public key
	registeredPubKey, err := LoadRegisteredPublicKey("internal/sftp/pub_key")
	if err != nil {
		log.Fatal("Failed to load registered public key: ", err)
	}
	// ssh client config
	config := &ssh.ClientConfig{
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: HostKeyCb(registeredPubKey),
	}

	// connect to ssh server
	conn, err := ssh.Dial("tcp", sshAddr, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	defer conn.Close()

	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var user = "USER"
	var sftpBasePath = fmt.Sprintf("home/%s", user)
	err = client.MkdirAll(sftpBasePath)
	if err != nil {
		log.Fatal(err)
	}

	// leave your mark
	f, err := client.Create(fmt.Sprintf("%s/hello.txt", sftpBasePath))
	if err != nil {
		log.Fatal(err)
	}

	if _, err := f.Write([]byte("Hello world!")); err != nil {
		log.Fatal(err)
	}

	// check error
	err = client.MkdirAll(sftpBasePath + "/dir")
	if err != nil {
		log.Fatal(err)
	}

	f.Close()

	// check it's there
	fi, err := client.Lstat("hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	if fi == nil {
		log.Fatal("file not found")
	}
	client.Close()

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

// TODO: Move this function to a separate package
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

// TODO: Move this function to a separate package
func HostKeyCb(registeredKey ssh.PublicKey) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		if string(key.Marshal()) == string(registeredKey.Marshal()) {
			return nil
		}

		return fmt.Errorf("host key mismatch")
	}
}

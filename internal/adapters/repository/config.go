package repository

import (
	"errors"
	"os"
	"time"
)

var (
	ErrConfigNotSet              = errors.New("config not set")
	ErrNoEnvSet                  = errors.New("no env set")
	ErrNoDigitalOceanCredentials = errors.New("no digital ocean credentials set")
)

type Config struct {
	env          string
	addr         string
	idleTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
	DB           DB
	auth         AuthStruct
	SFTP         SFTPStruct
	digitalOcean digitalOcean
}

type SFTPStruct struct {
	addr           string
	port           int
	publicKeyName  string
	publicKeyPath  string
	privateKeyPath string
}

type DB struct {
	DSN          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

type AuthStruct struct {
	privateKeyPath string
	publicKeyPath  string
	secretKey      string
}

type digitalOcean struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	bucket          string
}

func NewConfig(env string) (*Config, error) {
	c := &Config{
		env:          env,
		addr:         ":4000",
		idleTimeout:  2 * time.Minute,
		readTimeout:  5 * time.Second,
		writeTimeout: 10 * time.Second,
		DB: DB{
			DSN: os.Getenv("BUHO_DB_DSN"),
		},
		auth: AuthStruct{
			privateKeyPath: "internal/keys/jwt/private.pem",
			publicKeyPath:  "internal/keys/jwt/public.pem",
		},
		SFTP: SFTPStruct{
			addr:           "localhost",
			port:           2022,
			publicKeyPath:  "internal/keys/sftp/public.pem",
			privateKeyPath: "internal/keys/sftp/private.pem",
		},
		digitalOcean: digitalOcean{
			endpoint:        "fra1.digitaloceanspaces.com",
			accessKeyID:     os.Getenv("DO_SPACES_ACCESS_KEY"),
			secretAccessKey: os.Getenv("DO_SPACES_SECRET_KEY"),
			bucket:          "mussol",
		},
	}

	if c.digitalOcean.accessKeyID == "" || c.digitalOcean.secretAccessKey == "" {
		//	return nil, ErrNoDigitalOceanCredentials
	}

	return c, nil
}

package main

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
	db           db
	auth         authStruct
	sftp         sftpStruct
	digitalOcean digitalOcean
}

type sftpStruct struct {
	addr           string
	port           int
	publicKeyName  string
	publicKeyPath  string
	privateKeyPath string
}

type db struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

type authStruct struct {
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

func newConfig(env string) (*Config, error) {
	c := &Config{
		env:          env,
		addr:         ":4000",
		idleTimeout:  2 * time.Minute,
		readTimeout:  5 * time.Second,
		writeTimeout: 10 * time.Second,
		db: db{
			dsn: os.Getenv("BUHO_DB_DSN"),
		},
		auth: authStruct{
			privateKeyPath: "internal/keys/jwt/private.pem",
			publicKeyPath:  "internal/keys/jwt/public.pem",
		},
		sftp: sftpStruct{
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
		return nil, ErrNoDigitalOceanCredentials
	}

	// TODO: verify the config isn't empty
	return c, nil
}

package main

import (
	"os"
	"time"
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

func newConfig() *Config {
	c := &Config{
		env:          "development",
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
			publicKeyName:  "public.pem",
			publicKeyPath:  "internal/keys/sftp/public.pem",
			privateKeyPath: "internal/keys/sftp/private.pem",
		},
	}

	// TODO: verify the config isn't empty
	return c
}

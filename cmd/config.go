// Description: This file contains the configuration for the application.
package main

import (
	"os"
	"time"
)

type config struct {
	app  appConfig
	db   dbConfig
	auth authConfig
	sfpt sftpConfig
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

type sftpConfig struct {
	keyLocation string
	keyName     string
	addr        string
	port        int
}
type appConfig struct {
	env  string
	addr string
}

type authConfig struct {
	secretKey string
}

func NewConfig() *config {
	var cfg config
	cfg.app.env = os.Getenv("ENV")

	switch cfg.app.env {
	case "development":
		cfg = newDevelopmentConfig()

	case "production":
		cfg = newProductionConfig()

	default:
		cfg = newDevelopmentConfig()
	}

	//cfg.env = "development"
	//cfg.db.dsn = os.Getenv("BUHO_DB_DSN")
	//cfg.auth.secretKey = os.Getenv("BUHO_AUTH_SECRET_KEY")
	//cfg.env = "development"
	//cfg.db.dsn = os.Getenv("BUHO_DB_DSN")
	//cfg.auth.secretKey = os.Getenv("BUHO_AUTH_SECRET_KEY")

	return &cfg
}

// TODO: Add more configurations
func newDevelopmentConfig() config {
	return config{
		app: appConfig{
			env:  "development",
			addr: ":4000",
		},
		db: dbConfig{
			dsn:          os.Getenv("BUHO_DB_DSN"),
			maxOpenConns: 25,
			maxIdleConns: 25,
			maxIdleTime:  time.Minute,
		},
		auth: authConfig{
			secretKey: os.Getenv("BUHO_AUTH_SECRET_KEY"),
		},
	}
}

func newProductionConfig() config {
	return config{}
}

package models

import (
	"time"
)

type Config struct {
	Server Server
	DB     DB
}

type Server struct {
	Env          string
	Addr         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DB struct {
	DSN          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

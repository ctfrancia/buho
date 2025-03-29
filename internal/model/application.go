package model

// Environment defines the running environment
type Environment string

const (
	// Development environment will log to stdout
	Development Environment = "development"

	// Production environment will log to configured destinations
	Production Environment = "production"
)

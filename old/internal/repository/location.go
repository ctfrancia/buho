package repository

import (
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
}

/*
type Location struct {
	// gorm.Model
	// ID uint `gorm:"primaryKey"`
		City       string
		Country    string
		Province   string
		Address    string
		PostalCode string
}
*/

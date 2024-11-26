package model

import (
	"gorm.io/gorm"
)

type Tournament struct {
	gorm.Model
	Name        string
	Description string
	Date        string
}

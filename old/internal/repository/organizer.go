package repository

import (
	"gorm.io/gorm"
)

type Organizer struct {
	gorm.Model
}

/*
type OrganizerModel struct {
	ID      uint
	Name    string
	Email   string
	Phone   string
	Website string
}
*/

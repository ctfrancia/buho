// This file represents the actions and table definition for a player.
// A player is defined as someone who will participate in a tournament.
// A player is always a singular person.
package repository

import (
	"gorm.io/gorm"
)

// Player is a struct that defines the model in the Database
type Player struct {
	gorm.Model `json:"-"`
	// UUID is the unique identifier of the player and all actions against a player will be done using this UUID
	UUID string `json:"uuid" gorm:"not null"`
	// FirstName is the first name of the player
	FirstName string `json:"first_name" gorm:"not null"`
	// LastName is the last name of the player
	LastName string `json:"last_name" gorm:"not null"`
	// Email is the email of the player
	Email string `json:"email" gorm:"not null"`
	// Phone is the phone number of the player
	Phone string `json:"phone,omitempty"`
	// Country is the country of the player
	Country string `json:"country,omitempty"`
	// City is the city of the player
	City string `json:"city;omitempty"`
	// FideID is the FIDE ID of the player
	FideID string `json:"fide_id,omitempty"`
	// LocalID is the local ID of the player in the local chess federation
	LocalID string `json:"local_id,omitempty"`
	// LocalWebsite is the website of the local chess federation
	LocalWebsite string `json:"local_website,omitempty"`
	// Club is the club of the player
	Club string `json:"club,omitempty"`
}

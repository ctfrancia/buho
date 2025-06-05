package entities

import (
	"time"
)

type Tournament struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	StartDate   time.Time
	EndDate     time.Time
	CreatorID   uint
	PosterURL   string
	Description string
	QRCode      string
}

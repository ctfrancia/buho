package model

import ()

type CreateTournamentRequest struct {
	Name      string `json:"name,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

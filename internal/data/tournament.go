package data

type Tournament struct {
	ID        int    `json:"id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Name      string `json:"name,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Poster    string `json:"poster,omitempty"`
	IsValid   bool   `json:"is_valid,omitempty"`
}

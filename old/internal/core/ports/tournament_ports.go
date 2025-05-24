package ports

// TournamentRepository is the interface for interacting with the Tournament DB
type TournamentRepository interface {
	CreateTournament(tournament any) error
	/*
		GetByUUID(uuid string) (repository.Tournament, error)
		UpdateByUUID(uuid string, t repository.Tournament) error
		RemoveTournamentPosterURL(uuid string) error
		GetAll() ([]repository.Tournament, error)
		GetByDateRange(startDate, endDate string) ([]repository.Tournament, error)
	*/
}

// TournamentService is the interface for interacting with the Tournament Service
type TournamentService interface {
	CreateTournament(tournament any) error
	/*
		GetByUUID(uuid string) (repository.Tournament, error)
		UpdateByUUID(uuid string, t repository.Tournament) error
		RemoveTournamentPosterURL(uuid string) error
		GetAll() ([]repository.Tournament, error)
		GetByDateRange(startDate, endDate string) ([]repository.Tournament, error)
	*/
}

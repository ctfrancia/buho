package ports

// UserService is the interface for interacting with the User Service
type UserService interface {
	CreateUser(user any) error
	/*
		GetByUUID(uuid string) (repository.Tournament, error)
		UpdateByUUID(uuid string, t repository.Tournament) error
		RemoveTournamentPosterURL(uuid string) error
		GetAll() ([]repository.Tournament, error)
		GetByDateRange(startDate, endDate string) ([]repository.Tournament, error)
	*/
}

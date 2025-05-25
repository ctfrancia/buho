package repository

type HealthCheckRepository struct {
}

func NewHealthCheckRepository() *HealthCheckRepository {
	return &HealthCheckRepository{}
}

func (r *HealthCheckRepository) Check() error {
	return nil
}

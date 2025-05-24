package ports

type HealthCheckService interface {
	Check() error
}

type HealthCheckAdapter interface {
	HealthCheck() any
}

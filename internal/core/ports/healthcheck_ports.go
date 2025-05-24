package ports

import "net/http"

type HealthCheckHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type HealthCheckService interface {
	GetInformation()
}

type HealthCheckAdapter interface {
	HealthCheck() error
}

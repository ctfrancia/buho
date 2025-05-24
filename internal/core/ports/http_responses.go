package ports

import (
	"net/http"
)

type HttpResponse interface {
	WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error
	BadRequestResponse(w http.ResponseWriter, r *http.Request, err error)
	LogError(r *http.Request, err error)
	ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any)
	ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error)
	InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request)
	InvalidCredentialsCustomResponse(w http.ResponseWriter, r *http.Request, message string)
	ConflictResponse(w http.ResponseWriter, r *http.Request)
}

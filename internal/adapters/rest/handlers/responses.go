package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ctfrancia/buho/internal/core/ports"
)

type envelope map[string]any

type HandlerResponse struct {
	logger ports.Logger
}

func NewHandlerResponse(l ports.Logger) *HandlerResponse {
	return &HandlerResponse{
		logger: l,
	}
}

func (hr *HandlerResponse) WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	// FIXME: MarshalIndent is used to format the JSON output to make it more human-readable.
	// Change to json.Marshal if you want to remove the indentation.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (hr *HandlerResponse) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	hr.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (hr *HandlerResponse) LogError(r *http.Request, err error) {
	/*
		var (
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		fields := []zap.Field{
			zap.String("method", method),
			zap.String("uri", uri),
		}
		app.logger.Error(err.Error(), fields...)
	*/
}

func (hr *HandlerResponse) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := hr.WriteJSON(w, status, env, nil)
	if err != nil {
		hr.LogError(r, err)
		w.WriteHeader(500)
	}
}

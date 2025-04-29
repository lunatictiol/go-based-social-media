package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func init() {
	Validator = validator.New(validator.WithRequiredStructEnabled())
}
func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	max_bytes := 1_048_578

	r.Body = http.MaxBytesReader(w, r.Body, int64(max_bytes))

	decorder := json.NewDecoder(r.Body)
	decorder.DisallowUnknownFields()

	return decorder.Decode(data)
}

func WriteJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return WriteJSON(w, status, &envelope{Error: message})
}

func (a *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envolope struct {
		Success bool `json:success`
		Data    any  `json:data`
	}

	return WriteJSON(w, status, &envolope{
		Success: true,
		Data:    data,
	})

}

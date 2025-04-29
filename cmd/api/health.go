package main

import (
	"net/http"
)

func (a *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     a.config.env,
		"version": version,
	}
	if err := a.jsonResponse(w, http.StatusOK, data); err != nil {
		a.WriteInternalServerError(w, r, err)
	}
}

package main

import (
	"net/http"
)

func (a *application) WriteInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJSONError(w, http.StatusInternalServerError, "something went wrong")
}

func (a *application) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (a *application) NotfoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (a *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Errorf("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONError(w, http.StatusConflict, err.Error())
}

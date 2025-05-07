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

func (a *application) unauthorisedResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Errorf("unautorised response", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONError(w, http.StatusUnauthorized, err.Error())
}

func (a *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Warnf("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (a *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	a.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	WriteJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}

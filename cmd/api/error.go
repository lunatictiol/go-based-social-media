package main

import (
	"log"
	"net/http"
)

func (a *application) WriteInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error\n Method : %s path : %s error: %s", r.Method, r.URL.Path, err)
	WriteJSONError(w, http.StatusInternalServerError, "something went wrong")
}

func (app *application) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request\n Method : %s path : %s error: %s", r.Method, r.URL.Path, err)

	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) NotfoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found\n Method : %s path : %s error: %s", r.Method, r.URL.Path, err)

	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflict\n Method : %s path : %s error: %s", r.Method, r.URL.Path, err.Error())

	WriteJSONError(w, http.StatusConflict, err.Error())
}

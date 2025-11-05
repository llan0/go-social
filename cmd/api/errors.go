package main

import (
	"fmt"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Printf("internal server error: %s path: %s, error: %s", r.Method, r.URL.Path, err) // TODO: Better logger
	writeJSONError(w, http.StatusInternalServerError, "server ran into an issue!")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Printf("bad request error: %s path: %s, error: %s", r.Method, r.URL.Path, err) // TODO: Better logger
	writeJSONError(w, http.StatusBadRequest, err.Error())                              // why returning err as is here??
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Printf("not found error: %s path: %s, error: %s", r.Method, r.URL.Path, err) // TODO: Better logger
	writeJSONError(w, http.StatusNotFound, "not found")
}
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Printf("edit conflict: %s path: %s, error: %s", r.Method, r.URL.Path, err) // TODO: Better logger
	writeJSONError(w, http.StatusConflict, "edit conflict")                        // 409
}

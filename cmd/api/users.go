package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/llan0/go-social/internal/store"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	urlParam := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(urlParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	// fetch user
	user, err := app.store.Users.GetByID(ctx, userID)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

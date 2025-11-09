package main

import (
	"net/http"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// fetch post
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(1)) // TODO: hardcoded user atm, need to get this from auth session
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}

}

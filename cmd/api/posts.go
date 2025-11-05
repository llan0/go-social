package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/llan0/go-social/internal/store"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload // why not use store.Post? - because we are not using all the fields of it

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	//validate post
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//create a post
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1, //TODO: change after Auth
	}
	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusCreated, post); err != nil {
		// are we not sending the partial payload post here to writeJSON? NO, the other fields will be auto filled by db
		app.internalServerError(w, r, err)
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	// fetch post from middleware
	post := app.getPostFromCtx(r)

	//fectch comments
	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// add comments to post - since comments is not stored in post table
	post.Comments = comments

	if err := jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	urlParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(urlParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	// delete post
	if err = app.store.Posts.Delete(ctx, postID); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=100"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	// fetch post from middleware
	post := app.getPostFromCtx(r)

	//parse post
	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//validate struct
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	//update post
	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		// check for version conflict (optmistic concurrency control)
		switch {
		case errors.Is(err, store.ErrEditConflict):
			app.editConflictResponse(w, r, err)
			return
		case errors.Is(err, sql.ErrNoRows): // if your store still surfaces sql.ErrNoRows
			app.editConflictResponse(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// middleware for fetching post and adding to context
func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlParam := chi.URLParam(r, "postID")
		postID, err := strconv.ParseInt(urlParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		// fetch post
		post, err := app.store.Posts.GetByID(ctx, postID)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		//create a new context and insert the post
		ctx = context.WithValue(ctx, postCtx, post)

		// return the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// get the post from context
func (app *application) getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}

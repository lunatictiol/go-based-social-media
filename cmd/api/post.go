package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/lunatictiol/go-based-social-media/internal/store"
)

type postPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (a *application) createPosthandler(w http.ResponseWriter, r *http.Request) {
	var payload postPayload

	if err := ReadJSON(w, r, &payload); err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}
	if err := Validator.Struct(payload); err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}

	post := &store.Post{
		Content: payload.Content,
		Title:   payload.Title,
		UserId:  1,
		Tags:    payload.Tags,
	}

	ctx := context.Background()
	if err := a.store.Posts.Create(ctx, post); err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}

	if err := WriteJSON(w, http.StatusOK, post); err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}
}

func (a *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}
	ctx := r.Context()
	post, err := a.store.Posts.GetPostByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			a.NotfoundResponse(w, r, err)
		default:
			a.WriteInternalServerError(w, r, err)
		}
		return
	}

	comments, err := a.store.Comments.GetByPostID(ctx, id)
	if err != nil {
		a.WriteInternalServerError(w, r, err)
		return

	}
	post.Comments = comments

	if err := WriteJSON(w, http.StatusOK, post); err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}

}

func (a *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}
	ctx := r.Context()

	err = a.store.Posts.DeletePostByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			a.NotfoundResponse(w, r, err)
		default:
			a.WriteInternalServerError(w, r, err)
		}
		return

	}
	if err := WriteJSON(w, http.StatusOK, "post deleted successfully"); err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}

}

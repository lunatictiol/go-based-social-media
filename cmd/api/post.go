package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/lunatictiol/go-based-social-media/internal/store"
)

type contextKey string

const postKey contextKey = "post"

type postPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type updatePostPayload struct {
	Title   *string `json:"title" validate:"required,max=100"`
	Content *string `json:"content" validate:"required,max=1000"`
}
type commentPayload struct {
	UserId  int    `json:"user_id" validate:"required"`
	PostId  int    `json:"post_id" validate:"required"`
	Content string `json:"content" validate:"required,max=1000"`
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

	if err := a.jsonResponse(w, http.StatusOK, post); err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}
}

func (a *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := a.getPostfromCtx(r)
	comments, err := a.store.Comments.GetByPostID(r.Context(), post.Id)
	if err != nil {
		a.WriteInternalServerError(w, r, err)
		return

	}
	post.Comments = comments
	if err := a.jsonResponse(w, http.StatusOK, post); err != nil {
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
	if err := a.jsonResponse(w, http.StatusOK, "post deleted successfully"); err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}

}

func (a *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := a.getPostfromCtx(r)
	var payload updatePostPayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}
	if err := Validator.Struct(payload); err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	err = a.store.Posts.UpdatePost(r.Context(), post)
	if err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}

	if err := a.jsonResponse(w, http.StatusOK, post); err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}

}

func (a *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		ctx = context.WithValue(ctx, postKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (a *application) getPostfromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postKey).(*store.Post)
	return post
}

func (a *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload commentPayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}
	if err := Validator.Struct(payload); err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}
	err = a.store.Comments.Create(r.Context(), &store.Comment{
		UserId:  int64(payload.UserId),
		PostId:  int64(payload.PostId),
		Content: payload.Content,
	})

	if err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}

	if err := a.jsonResponse(w, http.StatusOK, "comment added successfully"); err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}

}

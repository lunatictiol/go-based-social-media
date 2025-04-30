package main

import (
	"net/http"

	"github.com/lunatictiol/go-based-social-media/internal/store"
)

func (a *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}
	fq, err := fq.Parse(r)
	if err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}
	if err := Validator.Struct(fq); err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}

	feed, err := a.store.Posts.GetUserFeed(ctx, int64(64), fq)
	if err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}
	if err = a.jsonResponse(w, http.StatusOK, feed); err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}
}

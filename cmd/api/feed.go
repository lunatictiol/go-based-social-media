package main

import "net/http"

func (a *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	feed, err := a.store.Posts.GetUserFeed(ctx, int64(64))
	if err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}
	if err = a.jsonResponse(w, http.StatusOK, feed); err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}
}

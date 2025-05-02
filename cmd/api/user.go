package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/lunatictiol/go-based-social-media/internal/store"
)

const userKey contextKey = "user"

func (a *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			a.WriteInternalServerError(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := a.store.Users.GetUserByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				a.NotfoundResponse(w, r, err)
			default:
				a.WriteInternalServerError(w, r, err)
			}
			return
		}
		ctx = context.WithValue(ctx, userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
func (a *application) getUserfromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userKey).(*store.User)
	return user
}

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/user/{userID} [get]
func (a *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := a.getUserfromCtx(r)
	if err := a.jsonResponse(w, http.StatusOK, user); err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}

}

type follower struct {
	userid int64 `json:"user_id" validation:"required"`
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/user/{userID}/follow [put]
func (a *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followeeuser := a.getUserfromCtx(r)

	var payload follower
	err := ReadJSON(w, r, &payload)
	if err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}
	if err := Validator.Struct(payload); err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}

	err = a.store.Followers.Follow(r.Context(), followeeuser.Id, payload.userid)
	if err != nil {
		switch err {
		case store.ErrConflict:
			a.conflictResponse(w, r, err)
			return
		default:
			a.WriteInternalServerError(w, r, err)
			return
		}
	}
	if err := a.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}

}

// UnfollowUser gdoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/user/{userID}/unfollow [put]
func (a *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followeeuser := a.getUserfromCtx(r)

	var payload follower
	err := ReadJSON(w, r, &payload)
	if err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}
	if err := Validator.Struct(payload); err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}

	err = a.store.Followers.UnFollow(r.Context(), followeeuser.Id, payload.userid)
	if err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}
	if err := a.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}

}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/user/activate/{token} [put]
func (a *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	err := a.store.Users.Activate(r.Context(), token)

	if err != nil {
		switch err {
		case store.ErrNotFound:
			a.NotfoundResponse(w, r, err)
		default:
			a.WriteInternalServerError(w, r, err)
		}
		return
	}

	if err := a.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		a.WriteInternalServerError(w, r, err)
	}
}

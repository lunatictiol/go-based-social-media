package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/lunatictiol/go-based-social-media/internal/store"
)

const userKey contextKey = "user"

func getUserfromCtx(r *http.Request) *store.User {
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
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}

	user, err := a.getUser(r.Context(), userID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			a.NotfoundResponse(w, r, err)
			return
		default:
			a.WriteInternalServerError(w, r, err)
			return
		}
	}

	if err := a.jsonResponse(w, http.StatusOK, user); err != nil {
		a.WriteInternalServerError(w, r, err)
	}
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
	followeruser := getUserfromCtx(r)

	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}

	err = a.store.Followers.Follow(r.Context(), followedID, followeruser.Id)
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
	followeeuser := getUserfromCtx(r)

	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}

	err = a.store.Followers.UnFollow(r.Context(), followeeuser.Id, followedID)
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

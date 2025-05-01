package main

import (
	"net/http"

	"github.com/lunatictiol/go-based-social-media/internal/store"
)

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required, max=100"`
	Email    string `json:"email" validate:"required,email, max=255"`
	Password string `json:"password" validate:"required,min=8, max=70"`
}

func (a *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := ReadJSON(w, r, payload); err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}

	if err := Validator.Struct(payload); err != nil {
		a.BadRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}
	if err := user.Password.Set(payload.Password); err != nil {
		a.WriteInternalServerError(w, r, err)
		return
	}
}

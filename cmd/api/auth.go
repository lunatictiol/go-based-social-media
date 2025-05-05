package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/lunatictiol/go-based-social-media/internal/mailer"
	"github.com/lunatictiol/go-based-social-media/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=70"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

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
//	@Router			/authenticate/user [post]
func (a *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := ReadJSON(w, r, &payload); err != nil {
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
	ctx := r.Context()
	plainToken := uuid.New().String()

	// hash the token for storage but keep the plain token for email
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err := a.store.Users.CreateAndInvite(ctx, user, hashToken, a.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			a.BadRequestResponse(w, r, err)
			return
		case store.ErrDuplicateUsername:
			a.BadRequestResponse(w, r, err)
			return
		default:
			a.WriteInternalServerError(w, r, err)
			return
		}
	}
	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	//send email
	activationURL := fmt.Sprintf("%s/confirm/%s", a.config.frontendURL, plainToken)
	isProdEnv := a.config.env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	// send mail
	status, err := a.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		a.logger.Errorw("error sending welcome email", "error", err)

		// rollback user creation if email fails (SAGA pattern)
		if err := a.store.Users.Delete(ctx, user.Id); err != nil {
			a.logger.Errorw("error deleting user", "error", err)
		}

		a.WriteInternalServerError(w, r, err)
		return
	}
	a.logger.Infow("Email sent", "status code", status)
	if err := a.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		a.WriteInternalServerError(w, r, err)
	}
}

package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/zrotrasukha/jobman/internal/data"
	"github.com/zrotrasukha/jobman/internal/validator"
)

func (app *application) CreateAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePlainText(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.User.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no matching email address found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrResponse(w, r, err)
		}
	}

	matches, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	if !matches {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Token.New(user.Id, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelop{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}
}

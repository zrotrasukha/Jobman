package main

import (
	"errors"
	"net/http"

	"github.com/zrotrasukha/jobman/internal/data"
	"github.com/zrotrasukha/jobman/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	v := validator.New()
	if data.ValidatePlainText(v, input.Password); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = user.Password.Set(input.Password)
	app.logger.Info("password set", "password", user.Password.Hash)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.User.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		default:
			app.serverErrResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusCreated, envelop{"user": user}, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}
}

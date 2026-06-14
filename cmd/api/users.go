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

	templateData := map[string]any{"Name": user.Name}
	app.background(func() {
		app.logger.Info("attempting to send welcome email", "recipient", user.Email)
		err := app.mailer.Send(user.Email, "welcome.tmpl", templateData)
		if err != nil {
			app.logger.Error("unable to send welcome email", "error", err)
			return
		}
		app.logger.Info("welcome email sent successfully", "recipient", user.Email)
	})

	err = app.writeJSON(w, http.StatusCreated, envelop{"user": user}, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}
}

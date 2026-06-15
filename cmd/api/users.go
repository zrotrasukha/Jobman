package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

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

	token, err := app.models.Token.New(user.Id, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}
	fmt.Println("token created successfully:", token.Plaintext)

	app.background(func() {
		templateData := map[string]any{
			"Name":            user.Name,
			"UserID":          user.Id,
			"ActivationToken": token.Plaintext,
		}
		app.logger.Info("attempting to send registration email", "recipient", user.Email)
		err := app.mailer.Send(user.Email, "registerUser.tmpl", templateData)
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

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlainText string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateTokenPlainText(v, input.TokenPlainText); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.User.GetForToken(input.TokenPlainText, data.ScopeActivation)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
			return
		default:
			app.serverErrResponse(w, r, err)
			return
		}
	}

	user.Activated = true

	err = app.models.User.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
			return
		default:
			app.serverErrResponse(w, r, err)
			return
		}
	}

	err = app.models.Token.DeleteAllforUser(user.Id, data.ScopeActivation)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"user": user}, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}
}

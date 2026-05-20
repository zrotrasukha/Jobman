package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/zrotrasukha/jobman/internal/data"
	"github.com/zrotrasukha/jobman/internal/validator"
)

func (app *application) CreateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Company_name string      `json:"company_name"`
		RoleTitle    string      `json:"role_title"`
		AppliedAt    string      `json:"applied_at"`
		Status       data.Status `json:"status"`
		Notes        string      `json:"notes"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var application = &data.JobApplication{
		CompanyName: input.Company_name,
		RoleTitle:   input.RoleTitle,
		Status:      input.Status,
		Notes:       input.Notes,
	}
	v := validator.New()

	if input.AppliedAt != "" {
		t, err := time.Parse(time.RFC3339, input.AppliedAt)
		if err != nil {
			v.AddError("applied_at", "must be a valid RFC3339 date")
		} else {
			application.AppliedAt = &t
		}
	}

	if data.ValidateJobApplication(v, application); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Application.Insert(application)
	if err != nil {
		app.serverErrResponse(w, r)
		return
	}

	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("/v1/applications/%d", application.ID))

	err = app.writeJSON(w, http.StatusCreated, envelop{"application": application}, header)
	if err != nil {
		app.serverErrResponse(w, r)
	}
}

func (app *application) GetApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("all the applications are here")
}

func (app *application) GetApplicationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	application, err := app.models.Application.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrResponse(w, r)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"application": application}, nil)
	if err != nil {
		app.serverErrResponse(w, r)
	}

}

func (app *application) UpdateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "update application")
}

func (app *application) DeleteApplicationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "delete application")
}

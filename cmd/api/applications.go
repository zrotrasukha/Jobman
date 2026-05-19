package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zrotrasukha/jobman/internal/data"
)

func (app *application) CreateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Company_name string      `json:"company_name"`
		RoleTitle    string      `json:"role_title"`
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

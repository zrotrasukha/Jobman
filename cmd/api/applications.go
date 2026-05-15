package main

import (
	"fmt"
	"net/http"

	"github.com/zrotrasukha/GO---Job-Application-Manager/internal/data"
)

func (app *application) CreateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Company_name      string `json:"company_name"`
		RoleTitle         string `json:"role_title"`
		Status            string `json:"status"`
		AppliedAt         string `json:"applied_at"`
		UpdatedAt         string `json:"updated_at"`
		LastCommunication string `json:"last_communication"`
		Notes             string `json:"notes"`
	}

	err := ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var application = &data.Application{
		CompanyName:       input.Company_name,
		RoleTitle:         input.RoleTitle,
		Status:            input.Status,
		AppliedAt:         input.AppliedAt,
		UpdatedAt:         input.UpdatedAt,
		LastCommunication: input.LastCommunication,
		Notes:             input.Notes,
	}

	err = app.models.Application.Insert(application)
	if err != nil {
		app.serverErrResponse(w, r)
		return
	}

	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("/v1/applications/%d", application.ID))

	err = WriteJSON(w, http.StatusCreated, envelop{"application": application}, header)
	if err != nil {
		app.serverErrResponse(w, r)
	}
}

func (app *application) GetApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "get applications")
}

func (app *application) GetApplicationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "get application")
}

func (app *application) UpdateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "update application")
}

func (app *application) DeleteApplicationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "delete application")
}

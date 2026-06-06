package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/zrotrasukha/jobman/internal/data"
	"github.com/zrotrasukha/jobman/internal/validator"
)

// Handler for creating a new job application
func (app *application) CreateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Company_name string      `json:"company_name"`
		RoleTitle    string      `json:"role_title"`
		AppliedAt    string      `json:"applied_at"`
		InterviewAt  string      `json:"interview_at"`
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
			application.AppliedAt = t
		}
	}

	// if an interview date is provided, set the stale_after field to 5 days after the interview date. If no interview date is provided, set the stale_after field to 30 days after the applied date.
	if input.InterviewAt != "" {
		t, err := time.Parse(time.RFC3339, input.InterviewAt)
		if err != nil {
			v.AddError("interview_at", "must be a valid RFC3339 date")
		} else {
			application.InterviewAt = &t
			gracePeriod := application.InterviewAt.Add(5 * 24 * time.Hour)
			application.StaleAfter = &gracePeriod
		}
	} else {
		if !application.AppliedAt.IsZero() {
			gracePeriod := application.AppliedAt.Add(30 * 24 * time.Hour)
			application.StaleAfter = &gracePeriod
		}
	}

	if data.ValidateJobApplication(v, application); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Application.Insert(application)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("/v1/applications/%d", application.ID))

	err = app.writeJSON(w, http.StatusCreated, envelop{"application": application}, header)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}
}

// Handler for listing job applications with optional search and pagination
func (app *application) ListApplicationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Search string
		data.Filters
	}

	qs := r.URL.Query()

	input.Search = app.readString(qs, "search", "")
	input.Filters.Page = app.readInt(qs, "page", 1)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "company_name", "role_title", "applied_at", "status", "-id", "-company_name", "-applied_at", "-role_title", "-status"}

	v := validator.New()

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	applications, metadata, err := app.models.Application.GetAll(input.Search, input.Filters)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"applications": applications, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}

}

// Handler for retrieving a specific job application by ID
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
			app.serverErrResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"application": application}, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}

}

// Handler for updating an existing job application by ID
func (app *application) UpdateApplicationHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readParamID(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	app.logger.Info("Updating application with ID", "id", id)

	application, err := app.models.Application.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrResponse(w, r, err)
			return
		}
	}

	var input struct {
		Company_name      *string      `json:"company_name"`
		RoleTitle         *string      `json:"role_title"`
		AppliedAt         *string      `json:"applied_at"`
		Interview_at      *string      `json:"interview_at"`
		LastCommunication *string      `json:"last_communication"`
		Status            *data.Status `json:"status"`
		Notes             *string      `json:"notes"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Company_name != nil {
		application.CompanyName = *input.Company_name
	}

	if input.RoleTitle != nil {
		application.RoleTitle = *input.RoleTitle
	}

	if input.Notes != nil {
		application.Notes = *input.Notes
	}

	if input.Status != nil {
		application.Status = *input.Status
	}

	if input.AppliedAt != nil {
		t, err := time.Parse(time.RFC3339, *input.AppliedAt)
		if err != nil {
			app.badRequestResponse(w, r, fmt.Errorf("invalid applied_at value: %w", err))
			return
		}
		application.AppliedAt = t
	}

	if input.LastCommunication != nil {
		t, err := time.Parse(time.RFC3339, *input.LastCommunication)
		if err != nil {
			app.badRequestResponse(w, r, fmt.Errorf("invalid last_communication value: %w", err))
			return
		}
		application.LastCommunication = &t
	}

	if input.Interview_at != nil {
		t, err := time.Parse(time.RFC3339, *input.Interview_at)
		if err != nil {
			app.badRequestResponse(w, r, fmt.Errorf("invalid interview_at value: %w", err))
			return
		}
		application.InterviewAt = &t
	}

	if input.Interview_at != nil {
		gracePeriod := application.InterviewAt.Add(5 * 24 * time.Hour)
		application.StaleAfter = &gracePeriod
	}

	v := validator.New()

	if data.ValidateJobApplication(v, application); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Application.Update(application)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.editConflictResponse(w, r)
			return
		default:
			app.serverErrResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"application": application}, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}
}

// Handler for deleting a specific job application by ID
func (app *application) DeleteApplicationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Application.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"message": "application successfully deleted"}, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}

}

package mocks

import (
	"context"
	"time"

	"github.com/zrotrasukha/jobman/internal/data"
)

var FixedDate = time.Date(2026, 8, 12, 11, 45, 0, 0, time.UTC)

type MockJobApplicationModel struct {
	InsertFunc                func(jobApp *data.JobApplication) error
	GETFunc                   func(id int64) (*data.JobApplication, error)
	GetAllFunc                func(searchString string, filters data.Filters) ([]*data.JobApplication, *data.Metadata, error)
	UpdateFunc                func(jobApp *data.JobApplication) error
	DeleteFunc                func(id int64) error
	MarkStaleApplicationsFunc func(ctx context.Context) (int64, error)
}

func (m MockJobApplicationModel) Insert(jobApp *data.JobApplication) error {
	if m.InsertFunc != nil {
		return m.InsertFunc(jobApp)
	}
	jobApp.ID = 1
	jobApp.CompanyName = "Test Company"
	jobApp.UpdatedAt = FixedDate
	jobApp.AppliedAt = FixedDate
	jobApp.RoleTitle = "Test Role"
	jobApp.Status = data.StatusApplied
	jobApp.Notes = "Test Notes"
	jobApp.Version = 1

	return nil
}

func (m MockJobApplicationModel) Get(id int64) (*data.JobApplication, error) {
	if m.GETFunc != nil {
		return m.GETFunc(id)
	}
	if id == 1 {
		return &data.JobApplication{
			ID:          1,
			CompanyName: "Test Company",
			RoleTitle:   "Test Role",
			Status:      data.StatusApplied,
			AppliedAt:   FixedDate,
			UpdatedAt:   FixedDate,
			Notes:       "Test Notes",
			Version:     1,
		}, nil
	}

	return nil, data.ErrRecordNotFound
}

func (m MockJobApplicationModel) GetAll(searchString string, filters data.Filters) ([]*data.JobApplication, *data.Metadata, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(searchString, filters)
	}
	jobApps := []*data.JobApplication{
		{
			ID:          1,
			CompanyName: "Test Company 1",
			RoleTitle:   "Test Role 1",
			Status:      data.StatusApplied,
			AppliedAt:   FixedDate,
			UpdatedAt:   FixedDate,
			Notes:       "Test Notes 1",
			Version:     1,
		},
		{
			ID:          2,
			CompanyName: "Test Company 2",
			RoleTitle:   "Test Role 2",
			Status:      data.StatusInterviewing,
			AppliedAt:   FixedDate,
			UpdatedAt:   FixedDate,
			Notes:       "Test Notes 2",
			Version:     1,
		},
	}
	metadata := &data.Metadata{
		CurrentPage:  1,
		PageSize:     20,
		FirstPage:    1,
		LastPage:     1,
		TotalRecords: 2,
	}
	return jobApps, metadata, nil
}

func (m MockJobApplicationModel) Update(jobApp *data.JobApplication) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(jobApp)
	}
	jobApp.UpdatedAt = FixedDate
	jobApp.Version++
	return nil
}

func (m MockJobApplicationModel) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	if id == 1 {
		return nil
	}
	return data.ErrRecordNotFound
}
func (m MockJobApplicationModel) MarkStaleApplications(ctx context.Context) (int64, error) {
	return 0, nil
}

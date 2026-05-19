package mocks

import "github.com/zrotrasukha/jobman/internal/data"

type MockJobApplicationModel struct {
	InsertFunc func(jobApp *data.JobApplication) error
	GETFunc    func(id int64) (*data.JobApplication, error)
}

func (m MockJobApplicationModel) Insert(jobApp *data.JobApplication) error {
	if m.InsertFunc != nil {
		return m.InsertFunc(jobApp)
	}
	jobApp.ID = 1
	jobApp.CompanyName = "Test Company"
	jobApp.RoleTitle = "Test Role"
	jobApp.Status = data.StatusApplied
	jobApp.Notes = "Test Notes"
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
			Notes:       "Test Notes",
		}, nil
	}

	return nil, data.ErrRecordNotFound
}

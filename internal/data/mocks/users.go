package mocks

import (
	"github.com/zrotrasukha/jobman/internal/data"
)

type MockUserModel struct {
	InsertFunc func(user *data.User) error
}

func (m MockUserModel) Insert(user *data.User) error {
	if m.InsertFunc != nil {
		return m.InsertFunc(user)
	}

	user.Id = 1
	user.Name = "Test User"
	user.Email = "test@user.com"
	user.CreatedAt = FixedDate
	return nil
}

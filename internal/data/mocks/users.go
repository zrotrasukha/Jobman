package mocks

import (
	"github.com/zrotrasukha/jobman/internal/data"
)

type MockUserModel struct {
	InsertFunc      func(user *data.User) error
	UpdateFunc      func(user *data.User) error
	GetForTokenFunc func(tokenScope, tokenPlaintext string) (*data.User, error)
	GetByEmailFunc  func(email string) (*data.User, error)
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

func (m MockUserModel) Update(user *data.User) error {
	return nil
}

func (m MockUserModel) GetForToken(tokenScope, tokenPlaintext string) (*data.User, error) {
	return &data.User{
		Id:        1,
		Name:      "Test User",
		Email:     "test@gmai.com",
		Activated: true,
		CreatedAt: FixedDate,
	}, nil
}

func (m MockUserModel) GetByEmail(email string) (*data.User, error) {
	return &data.User{
		Id:        1,
		Name:      "Test User",
		Email:     email,
		CreatedAt: FixedDate,
	}, nil
}

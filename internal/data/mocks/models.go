package mocks

import "github.com/zrotrasukha/jobman/internal/data"

func NewMockModels() data.Models {
	return data.Models{
		Application: MockJobApplicationModel{},
		User:        MockUserModel{},
	}
}

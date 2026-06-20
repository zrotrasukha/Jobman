package mocks

import (
	"time"

	"github.com/zrotrasukha/jobman/internal/data"
)

type MockTokenModel struct {
	NewFunc              func(userID int64, ttl time.Duration, scope string) (*data.Token, error)
	InsertFunc           func(token *data.Token) error
	DeleteAllforUserFunc func(userID int64, scope string) error
}

func (m MockTokenModel) New(userID int64, ttl time.Duration, scope string) (*data.Token, error) {
	if m.NewFunc != nil {
		return m.NewFunc(userID, ttl, scope)
	}
	return &data.Token{
		Plaintext: "3YBC5SUDHW2ADGZ6WD3VNDNV4I", // 26 chars
		UserID:    userID,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}, nil
}

func (m MockTokenModel) Insert(token *data.Token) error {
	if m.InsertFunc != nil {
		return m.InsertFunc(token)
	}
	return nil
}

func (m MockTokenModel) DeleteAllforUser(userID int64, scope string) error {
	if m.DeleteAllforUserFunc != nil {
		return m.DeleteAllforUserFunc(userID, scope)
	}
	return nil
}

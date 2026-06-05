package data_test

import (
	"errors"
	"testing"

	"github.com/zrotrasukha/jobman/internal/assert"
	"github.com/zrotrasukha/jobman/internal/data"
	testhelpers "github.com/zrotrasukha/jobman/internal/test_helpers"
)

func TestUserModel_Insert(t *testing.T) {
	pool := testhelpers.GetPool(t)
	testhelpers.TruncateTable(t, pool, testhelpers.TableUsers)

	models := data.NewModels(pool)

	t.Run("successful insert", func(t *testing.T) {
		user := &data.User{
			Name:  "Test User",
			Email: "test@gmail.com",
		}

		user.Password.Set("password")
		err := models.User.Insert(user)
		if err != nil {
			t.Fatal("unexpected error inserting user: ", err)
		}

		assert.Equal(t, user.Id, int64(1))
		assert.Equal(t, user.CreatedAt.IsZero(), false)
	})

	t.Run("duplicate email", func(t *testing.T) {
		testhelpers.TruncateTable(t, pool, testhelpers.TableUsers)

		users := []*data.User{
			{
				Name:  "Test User 1",
				Email: "test@gmail.com",
			},
			{
				Name:  "Test User 2",
				Email: "test@gmail.com",
			},
		}

		users[0].Password.Set("password")
		err := models.User.Insert(users[0])
		if err != nil {
			t.Fatal("unexpected error inserting user: ", err)
		}

		users[1].Password.Set("password")
		err = models.User.Insert(users[1])
		if !errors.Is(err, data.ErrDuplicateEmail) {
			t.Fatal("expected duplicate email error, got: ", err)
		}
	})

	t.Run("Duplicate email case insensitive", func(t *testing.T) {
		testhelpers.TruncateTable(t, pool, testhelpers.TableUsers)

		users := []*data.User{
			{
				Name:  "Test User 1",
				Email: "test@gmail.com",
			},
			{
				Name:  "Test User 2",
				Email: "TEST@GMAIL.COM",
			},
		}

		users[0].Password.Set("password")
		err := models.User.Insert(users[0])
		if err != nil {
			t.Fatal("unexpected error inserting user: ", err)
		}

		users[1].Password.Set("password")
		err = models.User.Insert(users[1])
		if !errors.Is(err, data.ErrDuplicateEmail) {
			t.Fatal("expected duplicate email error, got: ", err)
		}
	})
}

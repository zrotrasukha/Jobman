package data_test

import (
	"os"
	"testing"

	testhelpers "github.com/zrotrasukha/jobman/internal/test_helpers"
)

func TestMain(m *testing.M) {
	pool, terminate := testhelpers.NewTestPool()
	testhelpers.SharedPool = pool

	code := m.Run()

	terminate()
	os.Exit(code)
}

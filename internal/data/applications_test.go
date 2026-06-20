package data_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zrotrasukha/jobman/internal/assert"
	"github.com/zrotrasukha/jobman/internal/data"
	testhelpers "github.com/zrotrasukha/jobman/internal/test_helpers"
)

var date = time.Date(2026, 8, 12, 0, 0, 0, 0, time.UTC)

func createTestUser(t *testing.T, pool *pgxpool.Pool) *data.User {
	t.Helper()
	model := data.NewModels(pool)
	user := &data.User{
		Name:  "Test User",
		Email: "testuser@example.com",
	}
	user.Password.Set("password")
	err := model.User.Insert(user)
	if err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}
	return user
}

func TestJobApplicationModel_Insert(t *testing.T) {
	pool := testhelpers.GetPool(t)
	testhelpers.TruncateTable(t, pool, testhelpers.TableUsers)
	testhelpers.TruncateTable(t, pool, testhelpers.TableApplications)

	model := data.NewModels(pool)
	user := createTestUser(t, pool)

	inserted := &data.JobApplication{
		CompanyName: "Test Company",
		RoleTitle:   "Test Role",
		Status:      "Applied", // Casts cleanly to data.Status type if implicit, or use data.StatusApplied if defined
		AppliedAt:   date,
		Notes:       "This is a test job application.",
		UserID:      user.Id, // Changed from user.Id to user.ID
	}
	err := model.Application.Insert(inserted)
	if err != nil {
		t.Fatalf("inserted() returned an error: %v", err)
	}

	assert.Equal(t, inserted.ID, int64(1))
	assert.Equal(t, inserted.Version, int32(1))
	assert.Equal(t, inserted.AppliedAt.UTC(), date)
}

func TestJobApplicationModel_Get(t *testing.T) {
	pool := testhelpers.GetPool(t)
	model := data.NewModels(pool)

	testhelpers.TruncateTable(t, pool, testhelpers.TableUsers)
	testhelpers.TruncateTable(t, pool, testhelpers.TableApplications)
	user := createTestUser(t, pool)

	var now = time.Date(2026, 8, 12, 0, 0, 0, 0, time.UTC)

	t.Run("Existing ID", func(t *testing.T) {
		inserted := &data.JobApplication{
			CompanyName: "Test Company",
			RoleTitle:   "Test Role",
			Status:      "Applied",
			AppliedAt:   now,
			Notes:       "This is a test job application.",
			UserID:      user.Id, // Changed from user.Id to user.ID
		}

		if err := model.Application.Insert(inserted); err != nil {
			t.Fatalf("Error setting up insert %v", err)
		}

		got, err := model.Application.Get(inserted.ID, user.Id) // Changed from user.Id to user.ID
		if err != nil {
			t.Fatal("Error running model.Get()", err)
		}

		assert.Equal(t, got.ID, inserted.ID)
		assert.Equal(t, got.AppliedAt.Format(time.RFC3339), inserted.AppliedAt.Format(time.RFC3339))
		assert.Equal(t, got.UpdatedAt, inserted.UpdatedAt)
		assert.Equal(t, got.Version, inserted.Version)
	})

	t.Run("Non-existent ID", func(t *testing.T) {
		_, err := model.Application.Get(999999, user.Id)
		if !errors.Is(err, data.ErrRecordNotFound) {
			t.Errorf("want ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		_, err := model.Application.Get(-1, user.Id)
		if !errors.Is(err, data.ErrRecordNotFound) {
			t.Errorf("want ErrRecordNotFound, got %v", err)
		}
	})
}

func TestJobApplicationModel_Update(t *testing.T) {
	pool := testhelpers.GetPool(t)
	model := data.NewModels(pool)

	testhelpers.TruncateTable(t, pool, testhelpers.TableUsers)
	testhelpers.TruncateTable(t, pool, testhelpers.TableApplications)
	user := createTestUser(t, pool)

	t.Run("successful update", func(t *testing.T) {
		inserted := &data.JobApplication{
			CompanyName: "Test Company",
			RoleTitle:   "Test Role",
			Status:      "Applied",
			AppliedAt:   date,
			Notes:       "This is a test job application.",
			UserID:      user.Id,
		}

		err := model.Application.Insert(inserted)
		if err != nil {
			t.Fatalf("Error inserting: %v", err)
		}

		inserted.CompanyName = "Changed Company"

		if err = model.Application.Update(inserted, user.Id); err != nil {
			t.Fatalf("Error updating: %v", err)
		}

		got, err := model.Application.Get(inserted.ID, user.Id)
		if err != nil {
			t.Fatalf("Error fetching inserted job application: %v", err)
		}

		assert.Equal(t, got.Version, inserted.Version)
		assert.Equal(t, got.CompanyName, inserted.CompanyName)
	})

	t.Run("Optimistic Locking", func(t *testing.T) {
		testhelpers.TruncateTable(t, pool, testhelpers.TableApplications)

		inserted := &data.JobApplication{
			CompanyName: "Test Company",
			RoleTitle:   "Test Role",
			Status:      "Applied",
			AppliedAt:   date,
			Notes:       "This is a test job application.",
			UserID:      user.Id,
		}

		err := model.Application.Insert(inserted)
		if err != nil {
			t.Fatalf("Error inserting: %v", err)
		}

		v1, _ := model.Application.Get(inserted.ID, user.Id)
		v2, _ := model.Application.Get(inserted.ID, user.Id)

		v1.CompanyName = "Changed Company"

		if err = model.Application.Update(v1, user.Id); err != nil {
			t.Fatalf("Error updating v1: %v", err)
		}

		v2.RoleTitle = "Changed role"

		err = model.Application.Update(v2, user.Id)
		if !errors.Is(err, data.ErrEditConflict) {
			t.Fatalf("want: ErrEditConflict, go: %v ", err)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		testhelpers.TruncateTable(t, pool, testhelpers.TableApplications)

		inserted := &data.JobApplication{
			ID:          999999,
			CompanyName: "Test Company",
			RoleTitle:   "Test Role",
			Status:      "Applied",
			AppliedAt:   date,
			Notes:       "This is a test job application.",
			UserID:      user.Id,
		}

		err := model.Application.Update(inserted, user.Id)
		if !errors.Is(err, data.ErrEditConflict) {
			t.Fatalf("want: ErrEditConflict, go: %v ", err)
		}
	})
}

func TestJobApplicationModel_GetAll(t *testing.T) {
	pool := testhelpers.GetPool(t)
	model := data.NewModels(pool)

	testhelpers.TruncateTable(t, pool, testhelpers.TableUsers)
	testhelpers.TruncateTable(t, pool, testhelpers.TableApplications)
	user := createTestUser(t, pool)

	seeds := []*data.JobApplication{
		{CompanyName: "Google", RoleTitle: "Backend Engineer", Status: "Applied", AppliedAt: date, UserID: user.Id},
		{CompanyName: "Meta", RoleTitle: "Go Developer", Status: "Interviewing", AppliedAt: date, UserID: user.Id},
		{CompanyName: "Small Startup", RoleTitle: "Software Engineer", Status: "Rejected", AppliedAt: date, UserID: user.Id},
	}

	for _, ja := range seeds {
		err := model.Application.Insert(ja)
		if err != nil {
			t.Fatalf("Error inserting seeds %v", err)
		}
	}

	baseFilters := data.Filters{
		Page:         1,
		PageSize:     20,
		Sort:         "id",
		SortSafeList: []string{"id", "company_name", "role_title", "applied_at", "status", "-id", "-company_name", "-applied_at", "-role_title", "-status"},
	}

	t.Run("All results", func(t *testing.T) {
		applications, metadata, err := model.Application.GetAll("", baseFilters, user.Id)
		if err != nil {
			t.Fatalf("error running GetAll(): %v", err)
		}

		assert.Equal(t, len(applications), 3)
		assert.Equal(t, metadata.TotalRecords, 3)
	})

	t.Run("Company based full text search", func(t *testing.T) {
		applications, metadata, err := model.Application.GetAll("google", baseFilters, user.Id)
		if err != nil {
			t.Fatalf("Error running GetAll(): %v", err)
		}

		assert.Equal(t, metadata.TotalRecords, 1)
		assert.Equal(t, applications[0].CompanyName, seeds[0].CompanyName)
	})

	t.Run("Role based full text search", func(t *testing.T) {
		applications, metadata, err := model.Application.GetAll("Go Developer", baseFilters, user.Id)
		if err != nil {
			t.Fatalf("Error running GetAll(): %v", err)
		}

		assert.Equal(t, metadata.TotalRecords, 1)
		assert.Equal(t, applications[0].RoleTitle, "Go Developer")
	})

	t.Run("no results for unmatched search", func(t *testing.T) {
		_, metadata, err := model.Application.GetAll("thought police", baseFilters, user.Id)
		if err != nil {
			t.Fatalf("Error running GetAll(): %v", err)
		}

		assert.Equal(t, metadata.TotalRecords, 0)
	})

	t.Run("Pagination", func(t *testing.T) {
		f := data.Filters{
			Page:         1,
			PageSize:     2,
			Sort:         "id",
			SortSafeList: []string{"id"},
		}

		applications, metadata, err := model.Application.GetAll("", f, user.Id)
		if err != nil {
			t.Fatalf("error running GetAll(): %v", err)
		}

		assert.Equal(t, len(applications), 2)
		assert.Equal(t, metadata.TotalRecords, 3)
		assert.Equal(t, metadata.CurrentPage, 1)
	})

	t.Run("Second Page", func(t *testing.T) {
		f := data.Filters{
			Page:         2,
			PageSize:     2,
			Sort:         "id",
			SortSafeList: []string{"id"},
		}
		applications, metadata, err := model.Application.GetAll("", f, user.Id)
		if err != nil {
			t.Fatalf("error running GetAll(): %v", err)
		}

		assert.Equal(t, len(applications), 1)
		assert.Equal(t, metadata.TotalRecords, 3)
		assert.Equal(t, metadata.CurrentPage, 2)
	})

	t.Run("Sort by Company Name in ASC order", func(t *testing.T) {
		f := data.Filters{
			Page:         1,
			PageSize:     20,
			Sort:         "company_name",
			SortSafeList: []string{"id", "company_name"},
		}

		apps, _, err := model.Application.GetAll("", f, user.Id)
		if err != nil {
			t.Fatalf("error running GetAll(): %v", err)
		}

		for i, app := range apps {
			assert.Equal(t, app.CompanyName, seeds[i].CompanyName)
		}
	})

	t.Run("Sort by Company Name in DESC order", func(t *testing.T) {
		f := data.Filters{
			Page:         1,
			PageSize:     20,
			Sort:         "-company_name",
			SortSafeList: []string{"id", "-company_name"},
		}

		apps, _, err := model.Application.GetAll("", f, user.Id)
		if err != nil {
			t.Fatalf("error running GetAll(): %v", err)
		}

		for i := len(apps) - 1; i >= 0; i-- {
			assert.Equal(t, apps[len(apps)-1-i].CompanyName, seeds[i].CompanyName)
		}
	})
}

func TestJobApplicationModel_Delete(t *testing.T) {
	pool := testhelpers.GetPool(t)
	model := data.NewModels(pool)

	testhelpers.TruncateTable(t, pool, testhelpers.TableUsers)
	testhelpers.TruncateTable(t, pool, testhelpers.TableApplications)
	user := createTestUser(t, pool)

	t.Run("successful delete", func(t *testing.T) {
		testhelpers.TruncateTable(t, pool, testhelpers.TableApplications)

		inserted := &data.JobApplication{
			CompanyName: "Test Company",
			RoleTitle:   "Test Role",
			Status:      "Applied",
			AppliedAt:   date,
			Notes:       "This is a test job application.",
			UserID:      user.Id,
		}

		err := model.Application.Insert(inserted)
		if err != nil {
			t.Fatalf("Error inserting: %v", err)
		}

		err = model.Application.Delete(inserted.ID, user.Id)
		if err != nil {
			t.Fatalf("Error deleting: %v", err)
		}

		_, err = model.Application.Get(inserted.ID, user.Id)
		if !errors.Is(err, data.ErrRecordNotFound) {
			t.Errorf("want ErrRecordNotFound after delete, got %v", err)
		}
	})

	t.Run("Invalid Id delete", func(t *testing.T) {
		err := model.Application.Delete(-1, user.Id)
		if !errors.Is(err, data.ErrRecordNotFound) {
			t.Errorf("want ErrRecordNotFound for invalid id, got %v", err)
		}
	})

	t.Run("Double delete", func(t *testing.T) {
		testhelpers.TruncateTable(t, pool, testhelpers.TableApplications)
		inserted := &data.JobApplication{
			CompanyName: "Test Company",
			RoleTitle:   "Test Role",
			Status:      "Applied",
			AppliedAt:   date,
			UserID:      user.Id,
		}
		model.Application.Insert(inserted)
		model.Application.Delete(inserted.ID, user.Id) // first delete

		err := model.Application.Delete(inserted.ID, user.Id) // second delete
		if !errors.Is(err, data.ErrRecordNotFound) {
			t.Errorf("want ErrRecordNotFound on double delete, got %v", err)
		}
	})
}

func TestJobApplicationModel_MarkStaleApplication(t *testing.T) {
	pool := testhelpers.GetPool(t)
	model := data.NewModels(pool)
	testhelpers.TruncateTable(t, pool, testhelpers.TableUsers)
	testhelpers.TruncateTable(t, pool, testhelpers.TableApplications)
	user := createTestUser(t, pool)

	pastStale := time.Now().Add(-time.Hour)
	futureStale := time.Now().Add(10 * 24 * time.Hour)

	t.Run("mark stale applications", func(t *testing.T) {
		testhelpers.TruncateTable(t, pool, testhelpers.TableApplications)

		seed := []*data.JobApplication{
			{
				CompanyName: "Google",
				RoleTitle:   "Backend Engineer",
				Status:      "Applied",
				AppliedAt:   date,
				StaleAfter:  &pastStale,
				UserID:      user.Id,
			},
			{
				CompanyName: "Meta",
				RoleTitle:   "Go Developer",
				Status:      "Applied",
				AppliedAt:   date,
				StaleAfter:  &futureStale,
				UserID:      user.Id,
			},
		}

		for _, ja := range seed {
			err := model.Application.Insert(ja)
			if err != nil {
				t.Fatalf("Error inserting seeds %v", err)
			}
		}

		rowsAffected, err := model.Application.MarkStaleApplications(context.Background())
		if err != nil {
			t.Fatalf("Error running MarkStaleApplications: %v", err)
		}

		assert.Equal(t, rowsAffected, int64(1))

		apps, _, err := model.Application.GetAll("", data.Filters{
			Page:         1,
			PageSize:     20,
			Sort:         "id",
			SortSafeList: []string{"id"},
		}, user.Id)
		if err != nil {
			t.Fatalf("Error running GetAll(): %v", err)
		}

		assert.Equal(t, apps[0].Status, data.Status("Ghosted")) // Handled custom type
		assert.Equal(t, apps[1].Status, data.Status("Applied"))
	})

	t.Run("terminal status not affected", func(t *testing.T) {
		testhelpers.TruncateTable(t, pool, testhelpers.TableApplications)

		seeds := []*data.JobApplication{
			{
				CompanyName: "Google",
				RoleTitle:   "Backend Engineer",
				Status:      data.StatusOffered,
				AppliedAt:   date,
				StaleAfter:  &pastStale,
				UserID:      user.Id,
			},
			{
				CompanyName: "Meta",
				RoleTitle:   "Go Developer",
				Status:      data.StatusRejected,
				AppliedAt:   date,
				StaleAfter:  &pastStale,
				UserID:      user.Id,
			},
		}

		for _, seed := range seeds {
			err := model.Application.Insert(seed)
			if err != nil {
				t.Fatalf("Error inserting seeds %v", err)
			}
		}

		rowsAffected, err := model.Application.MarkStaleApplications(context.Background())
		if err != nil {
			t.Fatalf("Error running MarkStaleApplications: %v", err)
		}

		assert.Equal(t, rowsAffected, int64(0))
	})

	t.Run("already ghosted not affected", func(t *testing.T) {
		testhelpers.TruncateTable(t, pool, testhelpers.TableApplications)
		seed := []*data.JobApplication{
			{
				CompanyName: "Google",
				RoleTitle:   "Backend Engineer",
				Status:      data.StatusGhosted,
				AppliedAt:   date,
				StaleAfter:  &pastStale,
				UserID:      user.Id,
			},
			{
				CompanyName: "Meta",
				RoleTitle:   "Go Developer",
				Status:      data.StatusGhosted,
				AppliedAt:   date,
				StaleAfter:  &pastStale,
				UserID:      user.Id,
			},
		}

		for _, seed := range seed {
			err := model.Application.Insert(seed)
			if err != nil {
				t.Fatalf("Error inserting seeds %v", err)
			}
		}

		rowsAffected, err := model.Application.MarkStaleApplications(context.Background())
		if err != nil {
			t.Fatalf("Error running MarkStaleApplications: %v", err)
		}

		assert.Equal(t, rowsAffected, int64(0))
	})

	t.Run("zero applications to update", func(t *testing.T) {
		testhelpers.TruncateTable(t, pool, testhelpers.TableApplications)

		rowsAffected, err := model.Application.MarkStaleApplications(context.Background())
		if err != nil {
			t.Fatalf("Error running MarkStaleApplications: %v", err)
		}

		assert.Equal(t, rowsAffected, int64(0))
	})
}

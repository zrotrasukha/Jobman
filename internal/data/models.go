package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrRecordNotFound is returned when a requested record could not be found in the database. This error is used to indicate that a query did not return any results, and can be used by the application to provide appropriate feedback to the user or take other actions based on the absence of the requested data.
// ErrEditConflict is returned when an update operation fails due to a version mismatch, indicating that another concurrent update has occurred. This error is used to implement optimistic locking, allowing the application to detect and handle conflicts that arise when multiple users or processes attempt to modify the same record simultaneously.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Models is a wrapper struct that contains instances of all the model types used in the application. This allows for easy access to the models throughout the application, and provides a centralized location for managing database interactions. Each field in the Models struct corresponds to a specific model type, such as JobApplicationModelInterface for job applications.
type Models struct {
	Application JobApplicationModelInterface
}

// NewModels initializes a new Models struct with the provided database connection pool. It creates an instance of the JobApplicationModel and assigns it to the Application field of the Models struct. This function is typically called during application startup to set up the models with the necessary database connections, allowing the rest of the application to interact with the database through these models.
func NewModels(pool *pgxpool.Pool) Models {
	return Models{
		Application: NewJobApplicationModel(pool),
	}
}

// NewJobApplicationModel creates a new instance of the JobApplicationModel struct, which provides methods for interacting with the job applications table in the database. It takes a pgxpool.Pool as an argument, which is used to manage database connections and execute SQL queries. The returned JobApplicationModel can be used to perform CRUD operations on job application records in the database.
func NewJobApplicationModel(pool *pgxpool.Pool) JobApplicationModel {
	return JobApplicationModel{pool: pool}
}

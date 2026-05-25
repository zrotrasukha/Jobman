package data

import (
	"slices"
	"strings"

	"github.com/zrotrasukha/jobman/internal/validator"
)

// Filters struct defines the structure for pagination and sorting parameters used in API requests. It includes fields for the current page number, page size, sort order, and a list of allowed sort values. The struct is used to validate incoming request parameters and to calculate metadata for paginated responses.
type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

// Metadata struct defines the structure for pagination metadata included in API responses. It contains fields for the current page number, page size, first page number, last page number, and total number of records. This metadata helps clients understand the pagination context of the response and navigate through paginated results effectively.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// ValidateFilters checks the validity of the pagination and sorting parameters in the Filters struct. It ensures that the page number and page size are greater than zero, that the page size does not exceed a specified maximum (in this case, 100), and that the sort value is one of the allowed values defined in the SortSafeList. If any of these checks fail, appropriate error messages are added to the validator's error collection.
func ValidateFilters(v *validator.Validator, f Filters) {
	v.CheckField(f.Page > 0, "page", "must be greater than zero")
	v.CheckField(f.PageSize > 0, "page_size", "must be greater than zero")
	v.CheckField(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.CheckField(validator.PermittedValues(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

// calculateMetadata computes the pagination metadata based on the total number of records, the current page, and the page size. It calculates the first page (which is always 1), the last page (based on the total records and page size), and includes the current page and page size in the returned Metadata struct. If there are no records, it returns an empty Metadata struct.
func calculateMetadata(totalRecords, currentPage, pageSize int) *Metadata {
	if totalRecords == 0 {
		return &Metadata{}
	}

	return &Metadata{
		CurrentPage:  currentPage,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}

// Limit returns the page size from the Filters struct, which is used to limit the number of records returned in a paginated response.
func (f *Filters) Limit() int {
	return f.PageSize
}

// Offset calculates the offset for the SQL query based on the current page and page size. It determines how many records to skip before starting to return results for the current page.
func (f *Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

// SortColumn checks if the sort value in the Filters struct is included in the SortSafeList. If it is, it returns the sort column name by removing any leading "-" character (which indicates descending order). If the sort value is not in the safe list, it returns an empty string, indicating that no sorting should be applied.
func (f *Filters) SortColumn() string {
	if slices.Contains(f.SortSafeList, f.Sort) {
		return strings.TrimPrefix(f.Sort, "-")
	}

	return ""
}

// SortDirection determines the sort direction based on the presence of a leading "-" character in the Sort field of the Filters struct. If the Sort value starts with "-", it returns "DESC" for descending order; otherwise, it returns "ASC" for ascending order.
func (f *Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

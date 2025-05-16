package repository

import (
	"fmt"
	"math"
	"reflect"

	"gorm.io/gorm"
)

// PaginationParams holds the parameters for pagination: the current page and the number of items per page.
type PaginationParams struct {
	Page     int // Current page number
	PageSize int // Number of items per page
}

// PaginatedResult holds the results of a paginated query, including the data and pagination details.
type PaginatedResult struct {
	Data        interface{} `json:"data"`         // The actual data fetched from the database
	TotalCount  int64       `json:"total_count"`  // Total number of records available in the database
	TotalPages  int         `json:"total_pages"`  // Total number of pages based on the page size
	CurrentPage int         `json:"current_page"` // Current page being requested
	PageSize    int         `json:"page_size"`    // Number of items per page
}

// Paginate performs pagination on a database query using GORM.
// It takes a GORM database connection, pagination parameters, and a pointer to the result slice.
func Paginate(db *gorm.DB, params PaginationParams, result interface{}) (PaginatedResult, error) {
	// Get the type of the items in the result slice (using reflection)
	modelType := reflect.TypeOf(result).Elem().Elem()

	// Create a count query to count the total records in the table corresponding to the model type
	countQuery := db.Model(reflect.New(modelType).Interface())

	// Variable to hold the total count of records
	var totalCount int64
	// Execute the count query and handle any errors
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return PaginatedResult{}, fmt.Errorf("failed to count records: %w", err) // Return error if counting fails
	}

	// Calculate the total number of pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(params.PageSize)))

	// Calculate the offset for the SQL query based on the current page
	offset := (params.Page - 1) * params.PageSize

	// Execute the query to fetch the paginated records
	if err := db.Offset(offset).Limit(params.PageSize).Find(result).Error; err != nil {
		return PaginatedResult{}, fmt.Errorf("failed to fetch records: %w", err) // Return error if fetching fails
	}

	// Return the paginated result including data and pagination details
	return PaginatedResult{
		Data:        result,          // The fetched records
		TotalCount:  totalCount,      // Total number of records
		TotalPages:  totalPages,      // Total number of pages
		CurrentPage: params.Page,     // Current page number
		PageSize:    params.PageSize, // Number of items per page
	}, nil
}

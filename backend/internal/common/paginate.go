package common

import (
	"errors"
	"math"
	"net/http"
	"strconv"
)

// Paginate handles the pagination logic directly from the request context.
func Paginate(total int, data []interface{}, currentPage, perPage int) any {
	// Validate perPage to avoid division by zero or negative values
	if perPage <= 0 {
		perPage = 10 // Default value if not provided
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPages == 0 {
		totalPages = 1 // Ensure at least one page
	}

	// Validate currentPage
	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > totalPages {
		currentPage = totalPages
	}

	// Calculate nextPage and prevPage
	nextPage := currentPage + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	prevPage := currentPage - 1
	if prevPage < 1 {
		prevPage = 1
	}

	// Ensure data is not nil (optional: let caller handle this)
	if data == nil {
		data = []interface{}{}
	}

	// Return pagination result
	return map[string]any{
		"perPage":     perPage,
		"currentPage": currentPage,
		"nextPage":    nextPage,
		"prevPage":    prevPage,
		"totalPages":  totalPages,
		"total":       total,
		"data":        data,
	}
}

func GenericPaginate(total int, data interface{}, currentPage, perPage int) interface{} {
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	if currentPage > totalPages {
		currentPage = totalPages
	}

	nextPage := currentPage + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	prevPage := currentPage - 1
	if prevPage < 1 {
		prevPage = 1
	}
	// Create pagination result
	return map[string]any{
		"perPage":     perPage,
		"currentPage": currentPage,
		"nextPage":    nextPage,
		"prevPage":    prevPage,
		"totalPages":  totalPages,
		"total":       total,
		"data":        data,
	}
}

// Create a centralized method for parsing pagination and search parameters
func ParsePaginationParams(r *http.Request) (page, pageSize int, searchQuery string, err error) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	searchQuery = r.URL.Query().Get("search")

	// Default values
	page = 1
	pageSize = 10

	if pageStr != "" {
		pageParsed, err := strconv.Atoi(pageStr)
		if err != nil || pageParsed <= 0 {
			return 0, 0, "", errors.New("invalid page number")
		}
		page = pageParsed
	}

	if pageSizeStr != "" {
		pageSizeParsed, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSizeParsed <= 0 {
			return 0, 0, "", errors.New("invalid page size")
		}
		pageSize = pageSizeParsed
	}

	return page, pageSize, searchQuery, nil
}

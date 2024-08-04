package repository

import (
	"fmt"
	"math"
	"reflect"

	"gorm.io/gorm"
)

type PaginationParams struct {
	Page     int
	PageSize int
}

type PaginatedResult struct {
	Data        interface{} `json:"data"`
	TotalCount  int64       `json:"total_count"`
	TotalPages  int         `json:"total_pages"`
	CurrentPage int         `json:"current_page"`
	PageSize    int         `json:"page_size"`
}

func Paginate(db *gorm.DB, params PaginationParams, result interface{}) (PaginatedResult, error) {
	modelType := reflect.TypeOf(result).Elem().Elem() // Get the type of the slice

	countQuery := db.Model(reflect.New(modelType).Interface())

	var totalCount int64
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return PaginatedResult{}, fmt.Errorf("failed to count records: %w", err)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(params.PageSize)))

	offset := (params.Page - 1) * params.PageSize
	if err := db.Offset(offset).Limit(params.PageSize).Find(result).Error; err != nil {
		return PaginatedResult{}, fmt.Errorf("failed to fetch records: %w", err)
	}

	return PaginatedResult{
		Data:        result,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: params.Page,
		PageSize:    params.PageSize,
	}, nil
}

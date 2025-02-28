package repository

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

func UpdatePartial[T any](db *gorm.DB, dest T, updates T) error {
	// Use reflection to update only non-zero fields
	destValue := reflect.ValueOf(dest).Elem()
	updatesValue := reflect.ValueOf(updates)

	for i := 0; i < updatesValue.NumField(); i++ {
		updateField := updatesValue.Field(i)
		destField := destValue.Field(i)

		// Skip zero/default values
		if reflect.DeepEqual(updateField.Interface(), reflect.Zero(updateField.Type()).Interface()) {
			continue
		}

		// Set the field
		destField.Set(updateField)
	}

	if err := db.Save(dest).Error; err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	return nil
}

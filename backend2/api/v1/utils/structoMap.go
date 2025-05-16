package utils

import (
	"email-marketing-service/api/v1/model"
	"reflect"
)

// Instead of a method on UserResponse, make it a standalone function
func ToMap(u model.UserResponse) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(u)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		// Use json tag if available, otherwise use field name
		tag := field.Tag.Get("json")
		if tag == "" {
			tag = field.Name
		}

		// Skip if json tag is "-"
		if tag == "-" {
			continue
		}

		result[tag] = value
	}

	return result
}

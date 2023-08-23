package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func ValidateData(v interface{}) error {
	var validate = validator.New()

	err := validate.Struct(v)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		// Construct a response with validation errors
		errorMap := make(map[string]string)
		for _, e := range validationErrors {
			errorMap[e.Field()] = e.Tag()
		}
		errorResponse := map[string]interface{}{"errors": errorMap}

		return fmt.Errorf("validation errors: %v", errorResponse)

	}

	return nil

}

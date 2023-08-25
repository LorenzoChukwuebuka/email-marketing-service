package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidateData(v interface{}) error {
	validate := validator.New()

	err := validate.Struct(v)

	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			// Handle unexpected validation error type
			return fmt.Errorf("unexpected validation error: %v", err)
		}

		// Construct an error message with validation errors
		var errorMsgs []string
		for _, e := range validationErrors {
			errorMsgs = append(errorMsgs, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
		}
		return fmt.Errorf("validation errors: %s", strings.Join(errorMsgs, "; "))
	}

	return nil // No validation errors
}

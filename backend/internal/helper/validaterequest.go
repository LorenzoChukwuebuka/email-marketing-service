package helper

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

		var errorMsgs []string
		for _, e := range validationErrors {
			// Build more descriptive messages per rule
			var msg string
			switch e.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", e.Field())
			case "email":
				msg = fmt.Sprintf("%s must be a valid email address", e.Field())
			case "min":
				msg = fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param())
			case "max":
				msg = fmt.Sprintf("%s must be at most %s characters long", e.Field(), e.Param())
			case "omitempty":
				// usually you don’t need to show omitempty errors
				continue
			default:
				msg = fmt.Sprintf("%s is invalid (%s)", e.Field(), e.Tag())
			}

			errorMsgs = append(errorMsgs, msg)
		}

		return fmt.Errorf("validation errors: %s", strings.Join(errorMsgs, "; "))
	}

	return nil
}

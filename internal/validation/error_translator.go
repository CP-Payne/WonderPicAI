package validation

import (
	"fmt"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func TranslateValidationErrors(err error) (map[string]string, string) {
	fieldErrors := make(map[string]string)
	var generalError string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()
			tag := fieldError.Tag()

			var errorMessage string

			switch tag {
			case "required":
				errorMessage = fmt.Sprintf("%s is required.", fieldName)
			case "alphanum":
				errorMessage = fmt.Sprintf("%s must contain only letters and numbers.", fieldName)
			case "email":
				errorMessage = "Please enter a valid email address."
			case "min":
				errorMessage = fmt.Sprintf("%s must be at least %s characters long.", fieldName, fieldError.Param())
			case "max":
				errorMessage = fmt.Sprintf("%s cannot be longer than %s characters.", fieldName, fieldError.Param())
			case "eqfield":
				// Note: we are currently assuming that eqfield validator will only be used for password and ConfirmPassword
				// For eqfield, fieldError.Field() --> field being compared,
				// fieldError.Param() --> field it's compared against.
				errorMessage = "Passwords do not match."
			case "gte":
				errorMessage = fmt.Sprintf("%s must be greater than %s", fieldName, fieldError.Param())
			case "passwordcomplexity":
				errorMessage = "Password must be at least 8 characters long and include an uppercase letter, lowercase letter, number, and special character."
			default:
				errorMessage = fmt.Sprintf("%s is invalid.", fieldName)
			}

			if existing, ok := fieldErrors[fieldName]; ok {
				fieldErrors[lowerFirst(fieldName)] = existing + " " + errorMessage
			} else {
				fieldErrors[lowerFirst(fieldName)] = errorMessage
			}
		}
	} else if err != nil {
		generalError = "An unexpected error occured during validation. Please try again."
	}
	return fieldErrors, generalError
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

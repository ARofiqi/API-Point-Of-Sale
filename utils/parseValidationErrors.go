package utils

import "github.com/go-playground/validator/v10"

func ParseValidationErrors(err error) map[string]string {
	errorDetails := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errorDetails[err.Field()] = "Field validation failed on the '" + err.Tag() + "' tag"
	}
	return errorDetails
}

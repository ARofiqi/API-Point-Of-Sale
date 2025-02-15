package utils

import "github.com/go-playground/validator/v10"

func ParseValidationErrors(err error) map[string]string {
	errorDetails := make(map[string]string)
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			switch fieldErr.Tag() {
			case "required":
				errorDetails[fieldErr.Field()] = "Tidak boleh kosong"
			case "email":
				errorDetails[fieldErr.Field()] = "Format email tidak valid"
			case "min":
				errorDetails[fieldErr.Field()] = "Minimal " + fieldErr.Param() + " karakter"
			}
		}
	}
	return errorDetails
}

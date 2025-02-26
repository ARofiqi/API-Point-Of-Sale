package utils

import "github.com/go-playground/validator/v10"

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			fieldName := e.Field() // Nama field dalam struct
			var errMsg string

			// Sesuaikan pesan error sesuai tag validasi
			switch e.Tag() {
			case "required":
				errMsg = "tidak boleh kosong"
			case "min":
				errMsg = "terlalu pendek"
			case "max":
				errMsg = "terlalu panjang"
			case "email":
				errMsg = "format email tidak valid"
			default:
				errMsg = "tidak valid"
			}

			errors[fieldName] = errMsg
		}
	}

	return errors
}

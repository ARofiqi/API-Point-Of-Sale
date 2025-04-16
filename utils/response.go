package utils

import (
	"aro-shop/dto"

	"github.com/labstack/echo/v4"
)

func Response(c echo.Context, statusCode int, message string, data interface{}, err error, errorDetails dto.ErrorDetails) error {
	var errorID string

	// if errorDetails != nil && len(errorDetails) > 0 {
	// 	errorID := generateErrorID()
	// 	LogError(c, errorID, message, err)
	// }

	if err != nil {
		errorID := generateErrorID()
		LogError(c, errorID, message, err)
	}

	return c.JSON(statusCode, dto.Response{
		Message: message,
		Data:    data,
		Errors:  errorDetails,
		ErrorID: errorID,
	})
}

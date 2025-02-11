package utils

import (
	"aro-shop/models"

	"github.com/labstack/echo/v4"
)

func Response(c echo.Context, statusCode int, message string, data interface{}, err error, errorDetails map[string]string) error {
	var errorID string
	if errorDetails != nil && len(errorDetails) > 0 {
		errorID = generateErrorID()
		LogError(c, errorID, message, err)
	}

	return c.JSON(statusCode, models.Response{
		Message: message,
		Data:    data,
		Errors:  errorDetails,
		ErrorID: errorID,
	})
}

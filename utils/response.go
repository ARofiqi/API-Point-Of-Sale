package utils

import (
	"aro-shop/models"

	"github.com/labstack/echo/v4"
)

func Response(c echo.Context, statusCode int, message string, data interface{}, err error, errorDetails map[string]string) error {
	errorID := generateErrorID()

	if err != nil {
		LogError(c, errorID, message, err)
	}

	return c.JSON(statusCode, models.Response{
		Data:    data,
		Message: message,
		Errors:  []map[string]string{errorDetails},
		ErrorID: errorID,
	})
}

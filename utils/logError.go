package utils

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// "encoding/json"

func LogError(c echo.Context, errorID string, message string, err error) {
	requestBody := map[string]interface{}{}

	// if c.Request().Body != nil {
	// 	decoder := json.NewDecoder(c.Request().Body)
	// 	decoder.Decode(&requestBody)
	// }

	logrus.WithFields(logrus.Fields{
		"time":      time.Now().Format(time.RFC3339),
		", path":    c.Path(),
		", method":  c.Request().Method,
		", errorID": errorID,
		", error":   err.Error(),
		", request": requestBody,
	}).Error(message)
}

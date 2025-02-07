package utils

import (
	"encoding/json"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LogError(c echo.Context, errorID string, message string, err error) {
	requestBody := map[string]interface{}{}

	if c.Request().Body != nil {
		decoder := json.NewDecoder(c.Request().Body)
		decoder.Decode(&requestBody)
	}

	logrus.WithFields(logrus.Fields{
		"\n\ntime":         time.Now().Format(time.RFC3339),
		"\npath":         c.Path(),
		"\nmethod":       c.Request().Method,
		"\nerrorID":      errorID,
		"\nerrorMessage": err.Error(),
		"\nrequestBody":  requestBody,
	}).Error(message)
}

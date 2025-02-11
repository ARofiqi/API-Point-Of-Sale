package utils

import (
	"encoding/json"
	"io"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LogError(c echo.Context, errorID string, message string, err error) {
	requestBody := map[string]interface{}{}

	if c.Request().Body != nil {
		bodyBytes, readErr := io.ReadAll(c.Request().Body)
		if readErr == nil && len(bodyBytes) > 0 {
			json.Unmarshal(bodyBytes, &requestBody)
		}
	}

	logrus.WithFields(logrus.Fields{
		"time":      time.Now().Format(time.RFC3339),
		"path":      c.Path(),
		"method":    c.Request().Method,
		"errorID":   errorID,
		"error":     err.Error(),
		"request":   requestBody,
	}).Error(message)
}

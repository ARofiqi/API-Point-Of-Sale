package utils

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatal("Tidak dapat membuka file log: ", err)
	}

	logrus.SetOutput(io.MultiWriter(os.Stdout, logFile))
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func LogError(c echo.Context, errorID string, message string, err error) {
	requestBody := map[string]interface{}{}

	if c.Request().Body != nil {
		bodyBytes, readErr := io.ReadAll(c.Request().Body)
		if readErr == nil && len(bodyBytes) > 0 {
			json.Unmarshal(bodyBytes, &requestBody)
		}
	}

	var requestData interface{} = nil
	if len(requestBody) > 0 {
		requestData = requestBody
	}

	logrus.WithFields(logrus.Fields{
		"time":        time.Now().Format(time.RFC3339),
		"path":        c.Path(),
		"method":      c.Request().Method,
		"status_code": c.Response().Status,
		"errorID":     errorID,
		"error":       err.Error(),
		"request":     requestData,
		"user_agent":  c.Request().UserAgent(),
		"remote_ip":   c.RealIP(),
	}).Error(message)
}

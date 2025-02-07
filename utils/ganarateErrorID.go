package utils

import "time"

func generateErrorID() string {
	errorID := "ERR-" + time.Now().Format("150405")
	return errorID
}

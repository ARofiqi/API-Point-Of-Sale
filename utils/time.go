package utils

import "time"

// GetCurrentTime mengembalikan waktu saat ini dalam format UTC
func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

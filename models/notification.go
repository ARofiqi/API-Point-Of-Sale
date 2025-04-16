package models

import "time"
type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	IsRead    bool      `json:"is_read"`
}

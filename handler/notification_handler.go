package handler

import (
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetNotifications(c echo.Context) error {
	var notifications []models.Notification
	if err := db.DB.Order("created_at DESC").Find(&notifications).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch notifications", nil, err, nil)
	}
	return utils.Response(c, http.StatusOK, "Notifications retrieved successfully", notifications, nil, nil)
}

func MarkNotificationAsRead(c echo.Context) error {
	notificationID := c.Param("id")

	if err := db.DB.Model(&models.Notification{}).
		Where("id = ?", notificationID).
		Update("is_read", true).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to mark notification as read", nil, err, nil)
	}

	return utils.Response(c, http.StatusOK, "Notification marked as read", nil, nil, nil)
}

func CreateNotification(message string) error {
	notification := models.Notification{
		Message:   message,
		IsRead:    false,
		CreatedAt: utils.GetCurrentTime(),
	}

	if err := db.DB.Create(&notification).Error; err != nil {
		return err
	}
	return nil
}

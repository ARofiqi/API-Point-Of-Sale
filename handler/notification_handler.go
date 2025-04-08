package handler

import (
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/utils"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetNotifications(c echo.Context) error {
	var notifications []models.Notification
	var total int64

	// Ambil query parameter page & limit dengan default
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Default limit 10
	}

	offset := (page - 1) * limit

	// Hitung total notifikasi untuk pagination
	if err := db.DB.Model(&models.Notification{}).Count(&total).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to count notifications", nil, err, nil)
	}

	// Ambil data dengan pagination
	if err := db.DB.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notifications).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch notifications", nil, err, nil)
	}

	// Struktur respons dengan pagination info
	responseData := map[string]interface{}{
		"notifications": notifications,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": int(math.Ceil(float64(total) / float64(limit))),
		},
	}

	return utils.Response(c, http.StatusOK, "Notifications retrieved successfully", responseData, nil, nil)
}

func GetNotificationById(c echo.Context) error {
	var notification models.Notification

	// Ambil ID dari parameter URL
	id := c.Param("id")

	// Cari notifikasi berdasarkan ID
	if err := db.DB.First(&notification, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.Response(c, http.StatusNotFound, "Notification not found", nil, nil, nil)
		}
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch notification", nil, err, nil)
	}

	return utils.Response(c, http.StatusOK, "Notification retrieved successfully", notification, nil, nil)
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

package handler

import (
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/utils"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetPaymentMethods(c echo.Context) error {
	var methods []models.PaymentMethod
	if err := db.DB.Find(&methods).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch payment methods", nil, err, nil)
	}

	return utils.Response(c, http.StatusOK, "Payment methods retrieved successfully", methods, nil, nil)
}

func GetPaymentMethod(c echo.Context) error {
	id := c.Param("id")
	var method models.PaymentMethod

	if err := db.DB.First(&method, id).Error; err != nil {
		return utils.Response(c, http.StatusNotFound, "Payment method not found", nil, err, nil)
	}

	return utils.Response(c, http.StatusOK, "Payment method retrieved successfully", method, nil, nil)
}

func CreatePaymentMethod(c echo.Context) error {
	var input models.PaymentMethod
	if err := c.Bind(&input); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, nil)
	}

	if err := validate.Struct(input); err != nil {
		errDetails := utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, err, errDetails)
	}

	var existing models.PaymentMethod
	if err := db.DB.Where("name = ?", input.Name).First(&existing).Error; err == nil {
		return utils.Response(c, http.StatusBadRequest, "Payment method already exists", nil, nil, nil)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.Response(c, http.StatusInternalServerError, "Database error", nil, err, nil)
	}

	if err := db.DB.Create(&input).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to create payment method", nil, err, nil)
	}

	return utils.Response(c, http.StatusCreated, "Payment method created successfully", input, nil, nil)
}

func UpdatePaymentMethod(c echo.Context) error {
	id := c.Param("id")
	var method models.PaymentMethod

	if err := db.DB.First(&method, id).Error; err != nil {
		return utils.Response(c, http.StatusNotFound, "Payment method not found", nil, err, nil)
	}

	var input models.PaymentMethod
	if err := c.Bind(&input); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, nil)
	}

	method.Name = input.Name

	if err := db.DB.Save(&method).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to update payment method", nil, err, nil)
	}

	return utils.Response(c, http.StatusOK, "Payment method updated successfully", method, nil, nil)
}

func DeletePaymentMethod(c echo.Context) error {
	id := c.Param("id")
	if err := db.DB.Delete(&models.PaymentMethod{}, id).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to delete payment method", nil, err, nil)
	}

	return utils.Response(c, http.StatusOK, "Payment method deleted successfully", nil, nil, nil)
}

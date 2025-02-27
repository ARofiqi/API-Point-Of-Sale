package handler

import (
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetCategories(c echo.Context) error {
	var categories []models.Category
	if err := db.DB.Preload("Products").Find(&categories).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch categories", nil, err, nil)
	}
	return utils.Response(c, http.StatusOK, "Categories retrieved successfully", categories, nil, nil)
}

func GetCategory(c echo.Context) error {
	id := c.Param("id")
	var category models.Category
	if err := db.DB.Preload("Products").First(&category, id).Error; err != nil {
		return utils.Response(c, http.StatusNotFound, "Category not found", nil, err, nil)
	}
	return utils.Response(c, http.StatusOK, "Category retrieved successfully", category, nil, nil)
}

func CreateCategory(c echo.Context) error {
	var category models.Category
	if err := c.Bind(&category); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, nil)
	}

	if err := validate.Struct(category); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, err, nil)
	}

	if err := db.DB.Create(&category).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to create category", nil, err, nil)
	}
	return utils.Response(c, http.StatusCreated, "Category created successfully", category, nil, nil)
}

func UpdateCategory(c echo.Context) error {
	id := c.Param("id")
	var category models.Category
	if err := db.DB.First(&category, id).Error; err != nil {
		return utils.Response(c, http.StatusNotFound, "Category not found", nil, err, nil)
	}

	if err := c.Bind(&category); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, nil)
	}

	if err := validate.Struct(category); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, err, nil)
	}

	if err := db.DB.Save(&category).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to update category", nil, err, nil)
	}
	return utils.Response(c, http.StatusOK, "Category updated successfully", category, nil, nil)
}

func DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Category{}, id).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to delete category", nil, err, nil)
	}
	return utils.Response(c, http.StatusOK, "Category deleted successfully", nil, nil, nil)
}

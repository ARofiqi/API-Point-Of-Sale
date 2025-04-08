package handler

import (
	"aro-shop/cache"
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/utils"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetCategories(c echo.Context) error {
	cacheKey := "categories"

	// Cek cache Redis
	cachedData, err := cache.GetCache(cacheKey)
	if err == nil {
		var categories []models.Category
		if err := json.Unmarshal([]byte(cachedData), &categories); err == nil {
			return utils.Response(c, http.StatusOK, "Categories retrieved from cache", categories, nil, nil)
		}
	}

	// Fetch dari database jika cache tidak tersedia
	var categories []models.Category
	if err := db.DB.Find(&categories).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch categories", nil, err, nil)
	}

	// Simpan ke Redis dengan waktu kadaluarsa 5 menit
	jsonData, _ := json.Marshal(categories)
	cache.SetCache(cacheKey, string(jsonData), 5*time.Minute)

	return utils.Response(c, http.StatusOK, "Categories retrieved successfully", categories, nil, nil)
}

func GetCategoriesById(c echo.Context) error {
	id := c.Param("id")
	cacheKey := "category:" + id

	// Cek cache Redis
	cachedData, err := cache.GetCache(cacheKey)
	if err == nil {
		var category models.Category
		if err := json.Unmarshal([]byte(cachedData), &category); err == nil {
			return utils.Response(c, http.StatusOK, "Category retrieved from cache", category, nil, nil)
		}
	}

	// Fetch dari database jika cache tidak tersedia
	var category models.Category
	if err := db.DB.First(&category, id).Error; err != nil {
		return utils.Response(c, http.StatusNotFound, "Category not found", nil, err, nil)
	}

	// Simpan ke Redis dengan waktu kadaluarsa 5 menit
	jsonData, _ := json.Marshal(category)
	cache.SetCache(cacheKey, string(jsonData), 5*time.Minute)

	return utils.Response(c, http.StatusOK, "Category retrieved successfully", category, nil, nil)
}

func CreateCategory(c echo.Context) error {
	var category models.Category
	if err := c.Bind(&category); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, nil)
	}

	if category.Name == "" {
		return utils.Response(c, http.StatusBadRequest, "Name is required", nil, nil, nil)
	}

	if err := db.DB.Create(&category).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to create category", nil, err, nil)
	}

	// Hapus cache kategori agar data terbaru bisa diambil
	cache.DeleteCache("categories")

	return utils.Response(c, http.StatusCreated, "Category created successfully", category, nil, nil)
}

func UpdateCategory(c echo.Context) error {
	id := c.Param("id")
	var existingCategory models.Category
	if err := db.DB.First(&existingCategory, id).Error; err != nil {
		return utils.Response(c, http.StatusNotFound, "Category not found", nil, err, nil)
	}

	var updateData struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&updateData); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, nil)
	}

	if updateData.Name == "" {
		return utils.Response(c, http.StatusBadRequest, "Name is required", nil, nil, nil)
	}

	existingCategory.Name = updateData.Name

	if err := db.DB.Save(&existingCategory).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to update category", nil, err, nil)
	}

	cache.DeleteCache("categories")

	return utils.Response(c, http.StatusOK, "Category updated successfully", existingCategory, nil, nil)
}

func DeleteCategory(c echo.Context) error {
	id := c.Param("id")

	if err := db.DB.Delete(&models.Category{}, id).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to delete category", nil, err, nil)
	}

	// Hapus cache kategori terkait agar data terbaru bisa diambil
	cache.DeleteCache("categories")

	return utils.Response(c, http.StatusOK, "Category deleted successfully", nil, nil, nil)
}

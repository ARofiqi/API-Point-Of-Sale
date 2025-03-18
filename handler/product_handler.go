package handler

import (
	"aro-shop/cache"
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	validate = validator.New()
)

func resetRedish() {
	cache.DeleteCache("products_list:*:*:*:*")
	cache.DeleteCache("product:*")
	cache.DeleteCache("categories_with_products")
}

func GetProducts(c echo.Context) error {
	var (
		products     []models.Product
		errorDetails = make(models.ErrorDetails)
		category     = c.QueryParam("category")
		search       = c.QueryParam("search")
		page, _      = strconv.Atoi(c.QueryParam("page"))
		limit, _     = strconv.Atoi(c.QueryParam("limit"))
		cacheKey     = fmt.Sprintf("products_list:%s:%s:%d:%d", category, search, page, limit)
	)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// 1️⃣ Cek apakah data ada di Redis
	cachedData, err := cache.GetCache(cacheKey)
	if err == nil {
		var cachedProducts []models.ProductResponse
		json.Unmarshal([]byte(cachedData), &cachedProducts)
		return utils.Response(c, http.StatusOK, "Products fetched from cache", cachedProducts, nil, nil)
	}

	// 2️⃣ Jika tidak ada di Redis, ambil dari database
	query := db.DB.Preload("Category")
	if category != "" {
		query = query.Where("category_id = ?", category)
	}
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		errorDetails["database"] = "Gagal mengambil data produk dari database"
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch products", nil, err, errorDetails)
	}

	// 3️⃣ Konversi ke format response yang diinginkan
	var productResponses []models.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, models.ProductResponse{
			ID:       product.ID,
			Name:     product.Name,
			Price:    product.Price,
			Category: product.Category.Name,
		})
	}

	// 4️⃣ Simpan hasil query ke Redis untuk cache selama 10 menit
	jsonData, _ := json.Marshal(productResponses)
	cache.SetCache(cacheKey, string(jsonData), 10*time.Minute)

	return utils.Response(c, http.StatusOK, "Products fetched successfully", productResponses, nil, nil)
}

func GetProductByID(c echo.Context) error {
	var (
		id       = c.Param("id")
		product  models.Product
		cacheKey = fmt.Sprintf("product:%s", id)
	)

	// 1️⃣ Cek apakah data ada di Redis
	cachedData, err := cache.GetCache(cacheKey)
	if err == nil {
		var cachedProduct models.ProductResponse
		if json.Unmarshal([]byte(cachedData), &cachedProduct) == nil {
			return utils.Response(c, http.StatusOK, "Product fetched from cache", cachedProduct, nil, nil)
		}
	}

	// 2️⃣ Jika tidak ada di Redis, ambil dari database
	if err := db.DB.Preload("Category").First(&product, id).Error; err != nil {
		return utils.Response(c, http.StatusNotFound, "Product not found", nil, err, nil)
	}

	// 3️⃣ Konversi ke format response yang diinginkan
	productResponse := models.ConvertToProductResponse(product)

	// 4️⃣ Simpan hasil query ke Redis untuk cache selama 10 menit
	jsonData, _ := json.Marshal(productResponse)
	cache.SetCache(cacheKey, string(jsonData), 10*time.Minute)

	return utils.Response(c, http.StatusOK, "Product fetched successfully", productResponse, nil, nil)
}

func GetCategoriesWithProducts(c echo.Context) error {
	var (
		categories []models.Category
		cacheKey   = "categories_with_products"
	)

	// 1️⃣ Cek apakah data ada di Redis
	cachedData, err := cache.GetCache(cacheKey)
	if err == nil {
		json.Unmarshal([]byte(cachedData), &categories)
		return utils.Response(c, http.StatusOK, "Categories fetched from cache", categories, nil, nil)
	}

	// 2️⃣ Jika tidak ada di Redis, ambil dari database
	if err := db.DB.Preload("Products").Find(&categories).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch categories", nil, err, nil)
	}

	// 3️⃣ Simpan ke Redis dengan TTL 10 menit
	jsonData, _ := json.Marshal(categories)
	cache.SetCache(cacheKey, string(jsonData), 10*time.Minute)

	return utils.Response(c, http.StatusOK, "Categories fetched successfully", categories, nil, nil)
}

func CreateProduct(c echo.Context) error {
	var (
		category     models.Category
		product      models.Product
		errorDetails = make(models.ErrorDetails)
	)

	if err := c.Bind(&product); err != nil {
		// errorDetails["binding"] = err.Error()
		for _, err := range err.(validator.ValidationErrors) {
			errorDetails[err.Field()] = "Field validation failed on the '" + err.Tag() + "' tag"
		}
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, errorDetails)
	}

	if err := validate.Struct(product); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorDetails[err.Field()] = "Field validation failed on the '" + err.Tag() + "' tag"
		}
		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, err, errorDetails)
	}

	if err := db.DB.First(&category, product.CategoryID).Error; err != nil {
		errorDetails["category"] = "Category not found"
		return utils.Response(c, http.StatusBadRequest, "Invalid category ID", nil, err, errorDetails)
	}

	if err := db.DB.Create(&product).Error; err != nil {
		errorDetails["database"] = err.Error()
		return utils.Response(c, http.StatusInternalServerError, "Failed to create product", nil, err, errorDetails)
	}

	resetRedish()

	return utils.Response(c, http.StatusCreated, "Product created successfully", product, nil, nil)
}

func UpdateProduct(c echo.Context) error {
	var (
		errorDetails = make(models.ErrorDetails)
		product      models.Product
		input        models.Product
		id           = c.Param("id")
	)

	if err := db.DB.First(&product, id).Error; err != nil {
		return utils.Response(c, http.StatusNotFound, "Product not found", nil, err, nil)
	}

	if err := c.Bind(&input); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorDetails[err.Field()] = "Field validation failed on the '" + err.Tag() + "' tag"
		}
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, errorDetails)
	}

	if input.CategoryID != 0 {
		var category models.Category
		if err := db.DB.First(&category, input.CategoryID).Error; err != nil {
			return utils.Response(c, http.StatusBadRequest, "Invalid category ID", nil, err, nil)
		}
	}

	if err := db.DB.Model(&product).Updates(input).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to update product", nil, err, nil)
	}

	resetRedish()

	return utils.Response(c, http.StatusOK, "Product updated successfully", product, nil, nil)
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Product{}, id).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to delete product", nil, err, nil)
	}

	resetRedish()

	return utils.Response(c, http.StatusOK, "Product deleted successfully", nil, nil, nil)
}

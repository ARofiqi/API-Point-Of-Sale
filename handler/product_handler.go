package handler

import (
	"aro-shop/cache"
	"aro-shop/db"
	"aro-shop/dto"
	"aro-shop/models"
	"aro-shop/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	validate           = validator.New()
	cachedDataProducts = []string{
		"products_list:*",
		"product:*",
		"categories_with_products",
	}
)

func GetProducts(c echo.Context) error {
	var (
		products     []models.Product
		errorDetails = make(dto.ErrorDetails)
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

	// Cek apakah data ada di Redis
	cachedData, err := cache.GetCache(cacheKey)
	if err == nil {
		var cachedProducts []dto.ProductResponse
		json.Unmarshal([]byte(cachedData), &cachedProducts)
		return utils.Response(c, http.StatusOK, "Products fetched from cache", cachedProducts, nil, nil)
	}

	// Jika tidak ada di Redis, ambil dari database
	query := db.DB.Preload("Category")
	if category != "" {
		query = query.Where("category_id = ?", category)
	}
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		errorDetails["database"] = "Gagal mengambil data produk dari database"
		return utils.Response(c, http.StatusInternalServerError, "Internal server error", nil, err, errorDetails)
	}

	// Konversi ke format response yang diinginkan
	var productResponses []dto.ProductResponse
	if len(products) > 0 {
		for _, product := range products {
			productResponses = append(productResponses, dto.ConvertToProductResponse(product))
		}
	}

	// Simpan hasil query ke Redis untuk cache selama 10 menit
	jsonData, _ := json.Marshal(productResponses)
	cache.SetCache(cacheKey, string(jsonData), 10*time.Minute)

	return utils.Response(c, http.StatusOK, "Products fetched successfully", productResponses, nil, nil)
}

func GetProductByID(c echo.Context) error {
	var (
		id           = c.Param("id")
		product      models.Product
		cacheKey     = fmt.Sprintf("product:%s", id)
		errorDetails = make(dto.ErrorDetails)
	)

	// Cek apakah data ada di Redis
	cachedData, err := cache.GetCache(cacheKey)
	if err == nil {
		var cachedProduct dto.ProductResponse
		if json.Unmarshal([]byte(cachedData), &cachedProduct) == nil {
			return utils.Response(c, http.StatusOK, "Product fetched from cache", cachedProduct, nil, nil)
		}
	}

	// lakukan pengecekan id
	uuidID, err := uuid.Parse(id)
	if err != nil {
		errorDetails["id"] = "Invalid UUID format"
		return utils.Response(c, http.StatusBadRequest, "Invalid ID format", nil, err, errorDetails)
	}

	// Jika tidak ada di Redis, ambil dari database
	if err := db.DB.Preload("Category").First(&product, "id = ?", uuidID).Error; err != nil {
		errorDetails["id"] = "Product not found"
		return utils.Response(c, http.StatusNotFound, "Client error", nil, err, errorDetails)
	}

	// Konversi ke format response yang diinginkan
	productResponse := dto.ConvertToProductResponse(product)

	// Simpan hasil query ke Redis untuk cache selama 10 menit
	jsonData, _ := json.Marshal(productResponse)
	cache.SetCache(cacheKey, string(jsonData), 10*time.Minute)

	return utils.Response(c, http.StatusOK, "Product fetched successfully", productResponse, nil, nil)
}

func GetCategoriesWithProducts(c echo.Context) error {
	var (
		products    []models.Product
		cacheKey    = "categories_with_products"
		categoryMap = make(map[string][]dto.ProductResponse)
	)

	// Ambil query parameter `page` dan `limit`, default `page=1` dan `limit=5`
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 5 // Default 5 produk per kategori
	}

	// Cek apakah data ada di Redis
	cachedData, err := cache.GetCache(cacheKey)
	if err == nil && cachedData != "" {
		var cachedResponse []map[string]interface{}
		if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err == nil {
			return utils.Response(c, http.StatusOK, "Categories fetched from cache", cachedResponse, nil, nil)
		}
	}

	// Ambil semua produk dari database
	if err := db.DB.Preload("Category").Find(&products).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch products", nil, err, nil)
	}

	// Kelompokkan produk berdasarkan kategori
	for _, product := range products {
		categoryName := product.Category.Name
		categoryMap[categoryName] = append(categoryMap[categoryName], dto.ConvertToProductResponse(product))
	}

	// Format hasil sesuai output yang diinginkan
	var response []map[string]interface{}
	categoryID := 1
	for categoryName, productList := range categoryMap {
		// Hitung offset berdasarkan halaman
		offset := (page - 1) * limit
		end := offset + limit

		// Pastikan index tidak melebihi jumlah produk dalam kategori
		if offset > len(productList) {
			offset = len(productList)
		}
		if end > len(productList) {
			end = len(productList)
		}

		// Ambil produk sesuai pagination
		paginatedProducts := productList[offset:end]

		response = append(response, map[string]interface{}{
			"id":       categoryID,
			"name":     categoryName,
			"products": paginatedProducts,
			"page":     page,
			"limit":    limit,
			"total":    len(productList),
		})
		categoryID++
	}

	// Simpan ke Redis dengan TTL 10 menit
	jsonData, err := json.Marshal(response)
	if err == nil {
		cache.SetCache(cacheKey, string(jsonData), 10*time.Minute)
	}

	return utils.Response(c, http.StatusOK, "Categories fetched successfully", response, nil, nil)
}

func CreateProduct(c echo.Context) error {
	var (
		category     models.Category
		req          dto.ProductRequest
		product      models.Product
		errorDetails = make(dto.ErrorDetails)
	)

	if err := c.Bind(&req); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, nil)
	}

	if err := validate.Struct(req); err != nil {
		errorDetails := make(map[string]string)
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				errorDetails[e.Field()] = "Field validation failed on the '" + e.Tag() + "' tag"
			}
		}
		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, err, errorDetails)
	}

	product = models.Product{
		Name:       req.Name,
		Price:      req.Price,
		CategoryID: req.CategoryID,
	}

	if err := db.DB.First(&category, product.CategoryID).Error; err != nil {
		errorDetails["category"] = "Category not found"
		return utils.Response(c, http.StatusBadRequest, "Invalid category ID", nil, err, errorDetails)
	}

	if err := db.DB.Create(&product).Error; err != nil {
		errorDetails["database"] = err.Error()
		return utils.Response(c, http.StatusInternalServerError, "Failed to create product", nil, err, errorDetails)
	}

	if err := db.DB.Preload("Category").First(&product, product.ID).Error; err != nil {
		errorDetails["database"] = err.Error()
		return utils.Response(c, http.StatusInternalServerError, "Failed to load product with category", nil, err, errorDetails)
	}

	// Konversi ke format response yang diinginkan
	var productResponses dto.ProductResponse = dto.ConvertToProductResponse(product)

	go cache.ResetRedisCache(cachedDataProducts...)

	return utils.Response(c, http.StatusCreated, "Product created successfully", productResponses, nil, nil)
}

func UpdateProduct(c echo.Context) error {
	var (
		errorDetails = make(dto.ErrorDetails)
		product      models.Product
		input        models.Product
		id           = c.Param("id")
	)

	uuidID, err := uuid.Parse(id)
	if err != nil {
		errorDetails["id"] = "Invalid UUID format"
		return utils.Response(c, http.StatusBadRequest, "Invalid ID format", nil, err, errorDetails)
	}

	if err := db.DB.First(&product, "id = ?", uuidID).Error; err != nil {
		errorDetails["id"] = "Product not found"
		return utils.Response(c, http.StatusNotFound, "Client error", nil, err, errorDetails)
	}

	if err := c.Bind(&input); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorDetails[err.Field()] = "Field validation failed on the '" + err.Tag() + "' tag"
		}
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, errorDetails)
	}

	if input.CategoryID != uuid.Nil {
		var category models.Category
		if err := db.DB.First(&category, input.CategoryID).Error; err != nil {
			errorDetails["category"] = "Category not found"
			return utils.Response(c, http.StatusBadRequest, "Invalid category ID", nil, err, errorDetails)
		}
	}

	if err := db.DB.Model(&product).Updates(input).Error; err != nil {
		errorDetails["database"] = "Failed to update product"
		return utils.Response(c, http.StatusInternalServerError, "Internal server error", nil, err, errorDetails)
	}

	if err := db.DB.Preload("Category").First(&product, product.ID).Error; err != nil {
		errorDetails["database"] = err.Error()
		return utils.Response(c, http.StatusInternalServerError, "Failed to load product with category", nil, err, errorDetails)
	}

	var productResponses dto.ProductResponse = dto.ConvertToProductResponse(product)

	go cache.ResetRedisCache(cachedDataProducts...)

	return utils.Response(c, http.StatusOK, "Product updated successfully", productResponses, nil, nil)
}

func DeleteProduct(c echo.Context) error {
	var (
		id          = c.Param("id")
		errorDetail = make(dto.ErrorDetails)
	)

	// melakukan pengecekan id = uuid
	uuidID, err := uuid.Parse(id)
	if err != nil {
		errorDetail["id"] = "Invalid UUID format"
		return utils.Response(c, http.StatusBadRequest, "Invalid ID format", nil, err, errorDetail)
	}

	// mencari data di database
	if err := db.DB.Delete(&models.Product{}, "id = ?", uuidID).Error; err != nil {
		errorDetail["database"] = "Failed to delete product"
		return utils.Response(c, http.StatusInternalServerError, "Internal server error", nil, err, errorDetail)
	}

	// reset redis agar terjadi konsistensi data
	go cache.ResetRedisCache(cachedDataProducts...)

	return utils.Response(c, http.StatusOK, "Product deleted successfully", nil, nil, nil)
}

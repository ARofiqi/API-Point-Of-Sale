package handler

import (
	"aro-shop/cache"
	"aro-shop/db"
	"aro-shop/dto"
	"aro-shop/models"
	"aro-shop/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

	// Bind form field ke struct (bukan untuk file)
	req.Name = c.FormValue("name")
	req.Description = c.FormValue("description")
	req.Price, _ = strconv.ParseFloat(c.FormValue("price"), 64)
	req.Stock, _ = strconv.Atoi(c.FormValue("stock"))
	req.CategoryID, _ = uuid.Parse(c.FormValue("category_id"))

	// Ambil file dari form-data
	file, err := c.FormFile("image")
	if err != nil {
		errorDetails["image"] = "Image file is required"
		return utils.Response(c, http.StatusBadRequest, "Image is required", nil, err, errorDetails)
	}

	src, err := file.Open()
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to open image", nil, err, nil)
	}
	defer src.Close()

	// Simpan file di folder local
	filename := fmt.Sprintf("public/uploads/%s-%s", uuid.New().String(), file.Filename)
	dst, err := os.Create(filename)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to save image", nil, err, nil)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to copy image", nil, err, nil)
	}

	// Misal URL-nya dari server lokal
	baseURL := cfg.BaseURL
	imageURL := baseURL + "public/uploads/" + filepath.Base(filename)

	// Validasi manual jika perlu
	if err := validate.Struct(req); err != nil {
		errorDetails := make(map[string]string)
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				errorDetails[e.Field()] = "Field validation failed on the '" + e.Tag() + "' tag"
			}
		}
		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, err, errorDetails)
	}

	// Buat product
	product = models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		URLImage:    imageURL,
		CategoryID:  req.CategoryID,
	}

	// Cek kategori
	if err := db.DB.First(&category, product.CategoryID).Error; err != nil {
		errorDetails["category"] = "Category not found"
		return utils.Response(c, http.StatusBadRequest, "Invalid category ID", nil, err, errorDetails)
	}

	// Simpan ke DB
	if err := db.DB.Create(&product).Error; err != nil {
		errorDetails["database"] = err.Error()
		return utils.Response(c, http.StatusInternalServerError, "Failed to create product", nil, err, errorDetails)
	}

	// Ambil data lengkap + relasi
	if err := db.DB.Preload("Category").First(&product, product.ID).Error; err != nil {
		errorDetails["database"] = err.Error()
		return utils.Response(c, http.StatusInternalServerError, "Failed to load product with category", nil, err, errorDetails)
	}

	response := dto.ConvertToProductResponse(product)
	go cache.ResetRedisCache(cachedDataProducts...)

	return utils.Response(c, http.StatusCreated, "Product created successfully", response, nil, nil)
}

func UpdateProduct(c echo.Context) error {
	var (
		errorDetails = make(dto.ErrorDetails)
		product      models.Product
		id           = c.Param("id")
	)

	uuidID, err := uuid.Parse(id)
	if err != nil {
		errorDetails["id"] = "Invalid UUID format"
		return utils.Response(c, http.StatusBadRequest, "Invalid ID format", nil, err, errorDetails)
	}

	// Cek produk ada atau tidak
	if err := db.DB.First(&product, "id = ?", uuidID).Error; err != nil {
		errorDetails["id"] = "Product not found"
		return utils.Response(c, http.StatusNotFound, "Client error", nil, err, errorDetails)
	}

	// Ambil dan cek form value satu per satu
	if name := c.FormValue("name"); name != "" {
		product.Name = name
	}

	if description := c.FormValue("description"); description != "" {
		product.Description = description
	}

	if priceStr := c.FormValue("price"); priceStr != "" {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil || price <= 0 {
			errorDetails["price"] = "Invalid price"
		} else {
			product.Price = price
		}
	}

	if stockStr := c.FormValue("stock"); stockStr != "" {
		stock, err := strconv.Atoi(stockStr)
		if err != nil || stock < 0 {
			errorDetails["stock"] = "Invalid stock"
		} else {
			product.Stock = stock
		}
	}

	if categoryIDStr := c.FormValue("category_id"); categoryIDStr != "" {
		categoryID, err := uuid.Parse(categoryIDStr)
		if err != nil {
			errorDetails["category_id"] = "Invalid category ID"
		} else {
			var category models.Category
			if err := db.DB.First(&category, categoryID).Error; err != nil {
				errorDetails["category"] = "Category not found"
			} else {
				product.CategoryID = categoryID
			}
		}
	}

	if len(errorDetails) > 0 {
		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, nil, errorDetails)
	}

	// Upload file jika ada
	file, err := c.FormFile("url_image")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			errorDetails["file"] = "Failed to open file"
			return utils.Response(c, http.StatusInternalServerError, "File error", nil, err, errorDetails)
		}
		defer src.Close()

		filename := uuid.New().String() + filepath.Ext(file.Filename)
		filepath := "uploads/" + filename
		dst, err := os.Create(filepath)
		if err != nil {
			errorDetails["file"] = "Failed to save file"
			return utils.Response(c, http.StatusInternalServerError, "File error", nil, err, errorDetails)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			errorDetails["file"] = "Failed to copy file"
			return utils.Response(c, http.StatusInternalServerError, "File error", nil, err, errorDetails)
		}

		product.URLImage = c.Scheme() + "://" + c.Request().Host + "/uploads/" + filename
	}

	// Simpan perubahan
	if err := db.DB.Save(&product).Error; err != nil {
		errorDetails["database"] = "Failed to update product"
		return utils.Response(c, http.StatusInternalServerError, "Internal server error", nil, err, errorDetails)
	}

	if err := db.DB.Preload("Category").First(&product, product.ID).Error; err != nil {
		errorDetails["database"] = err.Error()
		return utils.Response(c, http.StatusInternalServerError, "Failed to load product with category", nil, err, errorDetails)
	}

	productResponses := dto.ConvertToProductResponse(product)
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

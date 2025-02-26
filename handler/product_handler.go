package handler

import (
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

func GetProducts(c echo.Context) error {
	category := c.QueryParam("category")
	search := c.QueryParam("search")

	var products []models.Product

	query := db.DB.Preload("Category")

	if category != "" {
		query = query.Where("category_id = ?", category)
	}

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Find(&products).Error; err != nil {
		errorDetails := map[string]string{"database": "Gagal mengambil data produk dari database"}
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch products", []models.Product{}, err, errorDetails)
	}

	if len(products) == 0 {
		return utils.Response(c, http.StatusOK, "No products found", []models.Product{}, nil, nil)
	}

	return utils.Response(c, http.StatusOK, "Products fetched successfully", products, nil, nil)
}

func CreateProduct(c echo.Context) error {
	var product models.Product
	if err := c.Bind(&product); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, nil)
	}

	if err := validate.Struct(product); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, err, nil)
	}

	var category models.Category
	if err := db.DB.First(&category, product.CategoryID).Error; err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid category ID", nil, err, nil)
	}

	if err := db.DB.Create(&product).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to create product", nil, err, nil)
	}

	return utils.Response(c, http.StatusCreated, "Product created successfully", product, nil, nil)
}

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product

	if err := db.DB.First(&product, id).Error; err != nil {
		return utils.Response(c, http.StatusNotFound, "Product not found", nil, err, nil)
	}

	var input models.Product
	if err := c.Bind(&input); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, nil)
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

	return utils.Response(c, http.StatusOK, "Product updated successfully", product, nil, nil)
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Product{}, id).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to delete product", nil, err, nil)
	}
	return utils.Response(c, http.StatusOK, "Product deleted successfully", nil, nil, nil)
}

func GetProductByID(c echo.Context) error {
	id := c.Param("id")
	var product models.Product

	if err := db.DB.Preload("Category").First(&product, id).Error; err != nil {
		return utils.Response(c, http.StatusNoContent, "Product not found", nil, err, nil)
	}

	return utils.Response(c, http.StatusOK, "Product fetched successfully", product, nil, nil)
}

func GetCategoriesWithProducts(c echo.Context) error {
	var categories []models.Category
	if err := db.DB.Preload("Products").Find(&categories).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch categories", nil, err, nil)
	}

	return utils.Response(c, http.StatusOK, "Categories with products fetched successfully", categories, nil, nil)
}

// package handler

// import (
// 	"aro-shop/db"
// 	"aro-shop/models"
// 	"aro-shop/utils"
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"time"

// 	"github.com/go-playground/validator/v10"
// 	"github.com/labstack/echo/v4"
// )

// var validate = validator.New()

// func GetProducts(c echo.Context) error {
// 	category := c.QueryParam("category")
// 	search := c.QueryParam("search")
// 	cacheKey := "products:all"
// 	ctx := context.Background()

// 	cachedProducts, err := db.RedisClient.Get(ctx, cacheKey).Result()
// 	if err == nil {
// 		var products []models.Product
// 		json.Unmarshal([]byte(cachedProducts), &products)
// 		return utils.Response(c, http.StatusOK, "Products fetched from cache", products, nil, nil)
// 	}

// 	var products []models.Product
// 	query := db.DB.Preload("Category")
// 	if category != "" {
// 		query = query.Where("category_id = ?", category)
// 	}
// 	if search != "" {
// 		query = query.Where("name LIKE ?", "%"+search+"%")
// 	}

// 	if err := query.Find(&products).Error; err != nil {
// 		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch products", []models.Product{}, err, nil)
// 	}

// 	if len(products) == 0 {
// 		return utils.Response(c, http.StatusOK, "No products found", []models.Product{}, nil, nil)
// 	}

// 	productsJSON, _ := json.Marshal(products)
// 	db.RedisClient.Set(ctx, cacheKey, productsJSON, 600*time.Second)

// 	return utils.Response(c, http.StatusOK, "Products fetched successfully", products, nil, nil)
// }

// func CreateProduct(c echo.Context) error {
// 	var product models.Product
// 	if err := c.Bind(&product); err != nil {
// 		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, nil)
// 	}

// 	if err := validate.Struct(product); err != nil {
// 		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, err, nil)
// 	}

// 	if err := db.DB.Create(&product).Error; err != nil {
// 		return utils.Response(c, http.StatusInternalServerError, "Failed to create product", nil, err, nil)
// 	}

// 	db.RedisClient.Del(context.Background(), "products:all")

// 	return utils.Response(c, http.StatusCreated, "Product created successfully", product, nil, nil)
// }

// func UpdateProduct(c echo.Context) error {
// 	id := c.Param("id")
// 	var product models.Product

// 	if err := db.DB.First(&product, id).Error; err != nil {
// 		return utils.Response(c, http.StatusNotFound, "Product not found", nil, err, nil)
// 	}

// 	var input models.Product
// 	if err := c.Bind(&input); err != nil {
// 		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, nil)
// 	}

// 	if err := db.DB.Model(&product).Updates(input).Error; err != nil {
// 		return utils.Response(c, http.StatusInternalServerError, "Failed to update product", nil, err, nil)
// 	}

// 	db.RedisClient.Del(context.Background(), "products:all", "product:"+id)

// 	return utils.Response(c, http.StatusOK, "Product updated successfully", product, nil, nil)
// }

// func DeleteProduct(c echo.Context) error {
// 	id := c.Param("id")
// 	if err := db.DB.Delete(&models.Product{}, id).Error; err != nil {
// 		return utils.Response(c, http.StatusInternalServerError, "Failed to delete product", nil, err, nil)
// 	}

// 	db.RedisClient.Del(context.Background(), "products:all", "product:"+id)

// 	return utils.Response(c, http.StatusOK, "Product deleted successfully", nil, nil, nil)
// }

// func GetProductByID(c echo.Context) error {
// 	id := c.Param("id")
// 	cacheKey := "product:" + id
// 	ctx := context.Background()

// 	cachedProduct, err := db.RedisClient.Get(ctx, cacheKey).Result()
// 	if err == nil {
// 		var product models.Product
// 		json.Unmarshal([]byte(cachedProduct), &product)
// 		return utils.Response(c, http.StatusOK, "Product fetched from cache", product, nil, nil)
// 	}

// 	var product models.Product
// 	if err := db.DB.Preload("Category").First(&product, id).Error; err != nil {
// 		return utils.Response(c, http.StatusNotFound, "Product not found", nil, err, nil)
// 	}

// 	productJSON, _ := json.Marshal(product)
// 	db.RedisClient.Set(ctx, cacheKey, productJSON, 600*time.Second)

// 	return utils.Response(c, http.StatusOK, "Product fetched successfully", product, nil, nil)
// }

// func GetCategoriesWithProducts(c echo.Context) error {
// 	cacheKey := "categories:all"
// 	ctx := context.Background()

// 	cachedCategories, err := db.RedisClient.Get(ctx, cacheKey).Result()
// 	if err == nil {
// 		var categories []models.Category
// 		json.Unmarshal([]byte(cachedCategories), &categories)
// 		return utils.Response(c, http.StatusOK, "Categories fetched from cache", categories, nil, nil)
// 	}

// 	var categories []models.Category
// 	if err := db.DB.Preload("Products").Find(&categories).Error; err != nil {
// 		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch categories", nil, err, nil)
// 	}

// 	categoriesJSON, _ := json.Marshal(categories)
// 	db.RedisClient.Set(ctx, cacheKey, categoriesJSON, 600*time.Second)

// 	return utils.Response(c, http.StatusOK, "Categories with products fetched successfully", categories, nil, nil)
// }

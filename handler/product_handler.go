package handler

import (
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/utils"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	validate     = validator.New()
	errorDetails = make(map[string]string)
)

func GetProducts(c echo.Context) error {
	category := c.QueryParam("category")
	search := c.QueryParam("search")

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
		SELECT products.id, products.name, products.price, products.category_id, category.name as category
		FROM products 
		INNER JOIN category ON products.category_id = category.id`)

	var args []interface{}
	var conditions []string

	if category != "" {
		conditions = append(conditions, "category_id = ?")
		args = append(args, category)
	}
	if search != "" {
		conditions = append(conditions, "products.name LIKE ?")
		args = append(args, "%"+search+"%")
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}

	rows, err := db.DB.Query(queryBuilder.String(), args...)
	if err != nil {
		utils.LogError(c, "ERR_FETCH_PRODUCTS", "Failed to fetch products", err)
		errorDetails["database"] = "Gagal mengambil data produk dari database"
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch products", []models.Product{}, err, errorDetails)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CategoryID, &p.Category); err != nil {
			utils.LogError(c, "ERR_PARSE_PRODUCTS", "Failed to parse products", err)
			errorDetails["parsing"] = "Gagal membaca data produk dari database"
			return utils.Response(c, http.StatusInternalServerError, "Failed to parse products", []models.Product{}, err, errorDetails)
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		utils.LogError(c, "ERR_ITERATE_ROWS", "Error iterating rows", err)
		errorDetails["database"] = "Terjadi kesalahan saat membaca hasil query"
		return utils.Response(c, http.StatusInternalServerError, "Error fetching products", []models.Product{}, err, errorDetails)
	}

	if len(products) == 0 {
		return utils.Response(c, http.StatusOK, "No products found", []models.Product{}, nil, nil)
	}

	return utils.Response(c, http.StatusOK, "Products fetched successfully", products, nil, nil)
}

func CreateProduct(c echo.Context) error {
	var product models.Product
	if err := c.Bind(&product); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, map[string]string{"request": "Format permintaan tidak valid"})
	}

	if err := validate.Struct(product); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errorDetails[e.Field()] = "Field " + e.Field() + " tidak valid atau kosong"
		}
		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, err, errorDetails)
	}

	stmt, err := db.DB.Prepare("INSERT INTO products (name, price, category_id) VALUES (?, ?, ?)")
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Database error", nil, err, nil)
	}
	defer stmt.Close()

	result, err := stmt.Exec(product.Name, product.Price, product.CategoryID)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to create product", nil, err, nil)
	}

	id, _ := result.LastInsertId()
	product.ID = int(id)
	return utils.Response(c, http.StatusCreated, "Product created successfully", product, nil, nil)
}

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := c.Bind(&product); err != nil {
		errorDetails["request"] = "Format permintaan tidak valid"
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, errorDetails)
	}

	var existingProduct models.Product
	err := db.DB.QueryRow("SELECT id FROM products WHERE id = ?", id).Scan(&existingProduct.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.Response(c, http.StatusNotFound, "Product with the given ID does not exist", nil, nil, nil)
		}
		return utils.Response(c, http.StatusInternalServerError, "Database error", nil, err, nil)
	}

	stmt, err := db.DB.Prepare("UPDATE products SET name = ?, price = ?, category_id = ? WHERE id = ?")
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Database error", nil, err, nil)
	}
	defer stmt.Close()

	result, err := stmt.Exec(product.Name, product.Price, product.CategoryID, id)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to update product", nil, err, nil)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return utils.Response(c, http.StatusOK, "No changes detected, product remains the same", product, nil, nil)
	}

	return utils.Response(c, http.StatusOK, "Product updated successfully", product, nil, nil)
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	stmt, err := db.DB.Prepare("DELETE FROM products WHERE id = ?")
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Database error", nil, err, nil)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to delete product", nil, err, nil)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return utils.Response(c, http.StatusNotFound, "Product not found", nil, nil, nil)
	}

	return utils.Response(c, http.StatusOK, "Product deleted successfully", nil, nil, nil)
}

func GetProductByID(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	id := c.Param("id")
	var product models.Product

	err := db.DB.QueryRowContext(ctx, `
		SELECT products.id, products.name, products.price, products.category_id, category.name as category 
		FROM products 
		INNER JOIN category ON products.category_id = category.id 
		WHERE products.id = ?`, id).
		Scan(&product.ID, &product.Name, &product.Price, &product.CategoryID, &product.Category)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			errorDetails["database"] = "Produk dengan ID tersebut tidak ditemukan"
			return utils.Response(c, http.StatusNotFound, "Product not found", nil, err, errorDetails)
		}

		errorDetails["database"] = "Terjadi kesalahan saat mengambil data produk"
		return utils.Response(c, http.StatusInternalServerError, "Internal server error", nil, err, errorDetails)
	}

	return utils.Response(c, http.StatusOK, "Product fetched successfully", product, nil, nil)
}

func GetCategoriesWithProducts(c echo.Context) error {
	query := `
		SELECT category.id, category.name AS category, products.id, products.name, products.price
		FROM products
		INNER JOIN category ON products.category_id = category.id
		ORDER BY category.id
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		errorDetails["database"] = "Gagal mengambil data kategori dengan produk"
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch categories", nil, err, errorDetails)
	}
	defer rows.Close()

	categories := make(map[string][]models.Product)
	for rows.Next() {
		var categoryID int
		var categoryName string
		var product models.Product

		if err := rows.Scan(&categoryID, &categoryName, &product.ID, &product.Name, &product.Price); err != nil {
			errorDetails["parsing"] = "Gagal membaca data kategori dari database"
			return utils.Response(c, http.StatusInternalServerError, "Failed to parse categories", nil, err, errorDetails)
		}

		categories[categoryName] = append(categories[categoryName], product)
	}

	return utils.Response(c, http.StatusOK, "Categories with products fetched successfully", categories, nil, nil)
}

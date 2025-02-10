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

	errorDetails := make(map[string]string)

	query := "SELECT id, name, price, category FROM products"
	var args []interface{}

	if category != "" {
		query += " WHERE category = ?"
		args = append(args, category)
	}

	if search != "" {
		if category != "" {
			query += " AND name LIKE ?"
		} else {
			query += " WHERE name LIKE ?"
		}
		args = append(args, "%"+search+"%")
	}

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		utils.LogError(c, "ERR_FETCH_PRODUCTS", "Failed to fetch products", err)
		errorDetails["database"] = "Gagal mengambil data produk dari database"
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch products", []models.Product{}, err, errorDetails)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category); err != nil {
			utils.LogError(c, "ERR_PARSE_PRODUCTS", "Failed to parse products", err)
			errorDetails["parsing"] = "Gagal membaca data produk dari database"
			return utils.Response(c, http.StatusInternalServerError, "Failed to parse products", []models.Product{}, err, errorDetails)
		}
		products = append(products, p)
	}

	return utils.Response(c, http.StatusOK, "Products fetched successfully", products, nil, errorDetails)
}

func CreateProduct(c echo.Context) error {
	errorDetails := make(map[string]string)

	var products models.Product
	if err := c.Bind(&products); err != nil {
		utils.LogError(c, "ERR_BIND_PRODUCT", "Invalid request format", err)
		errorDetails["request"] = "Format permintaan tidak valid"
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, errorDetails)
	}

	if err := validate.Struct(products); err != nil {
		utils.LogError(c, "ERR_VALIDATE_PRODUCT", "Validation failed", err)

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				errorDetails[e.Field()] = "Field " + e.Field() + " tidak valid atau kosong"
			}
		}

		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, err, errorDetails)
	}

	result, err := db.DB.Exec("INSERT INTO products (name, price, category) VALUES (?, ?, ?)", products.Name, products.Price, products.Category)
	if err != nil {
		utils.LogError(c, "ERR_CREATE_PRODUCT", "Failed to create product", err)
		errorDetails["database"] = "Gagal menyimpan produk ke database"
		return utils.Response(c, http.StatusInternalServerError, "Failed to create product", nil, err, errorDetails)
	}

	id, _ := result.LastInsertId()
	products.ID = int(id)
	return utils.Response(c, http.StatusCreated, "Product created successfully", products, nil, errorDetails)
}

func UpdateProduct(c echo.Context) error {
	errorDetails := make(map[string]string)

	id := c.Param("id")
	var p models.Product
	if err := c.Bind(&p); err != nil {
		utils.LogError(c, "ERR_BIND_UPDATE_PRODUCT", "Invalid request format", err)
		errorDetails["request"] = "Format permintaan tidak valid"
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, errorDetails)
	}

	_, err := db.DB.Exec("UPDATE products SET name = ?, price = ?, category = ? WHERE id = ?", p.Name, p.Price, p.Category, id)
	if err != nil {
		utils.LogError(c, "ERR_UPDATE_PRODUCT", "Failed to update product", err)
		errorDetails["database"] = "Gagal memperbarui produk dalam database"
		return utils.Response(c, http.StatusInternalServerError, "Failed to update product", nil, err, errorDetails)
	}

	return utils.Response(c, http.StatusOK, "Product updated successfully", p, nil, nil)
}

func DeleteProduct(c echo.Context) error {
	errorDetails := make(map[string]string)

	id := c.Param("id")
	result, err := db.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		utils.LogError(c, "ERR_DELETE_PRODUCT", "Failed to delete product", err)
		errorDetails["database"] = "Gagal menghapus produk dari database"
		return utils.Response(c, http.StatusInternalServerError, "Failed to delete product", nil, err, errorDetails)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		errorDetails["id"] = "Produk dengan ID tersebut tidak ditemukan"
		return utils.Response(c, http.StatusNotFound, "Product not found", nil, nil, errorDetails)
	}

	return utils.Response(c, http.StatusOK, "Product deleted successfully", nil, nil, nil)
}

func GetProductByID(c echo.Context) error {
	errorDetails := make(map[string]string)
	id := c.Param("id")

	var product models.Product
	err := db.DB.QueryRow("SELECT id, name, price, category FROM products WHERE id = ?", id).
		Scan(&product.ID, &product.Name, &product.Price, &product.Category)

	if err != nil {
		utils.LogError(c, "ERR_FETCH_PRODUCT_BY_ID", "Failed to fetch product", err)
		errorDetails["database"] = "Produk dengan ID tersebut tidak ditemukan"
		return utils.Response(c, http.StatusNotFound, "Product not found", nil, err, errorDetails)
	}

	return utils.Response(c, http.StatusOK, "Product fetched successfully", product, nil, nil)
}

func GetCategoriesWithProducts(c echo.Context) error {
	errorDetails := make(map[string]string)

	query := "SELECT category, id, name, price FROM products ORDER BY category"
	rows, err := db.DB.Query(query)
	if err != nil {
		utils.LogError(c, "ERR_FETCH_CATEGORIES", "Failed to fetch categories with products", err)
		errorDetails["database"] = "Gagal mengambil data kategori dengan produk"
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch categories", nil, err, errorDetails)
	}
	defer rows.Close()

	categories := make(map[string][]models.Product)
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.Category, &p.ID, &p.Name, &p.Price); err != nil {
			utils.LogError(c, "ERR_PARSE_CATEGORIES", "Failed to parse category data", err)
			errorDetails["parsing"] = "Gagal membaca data kategori dari database"
			return utils.Response(c, http.StatusInternalServerError, "Failed to parse categories", nil, err, errorDetails)
		}
		categories[p.Category] = append(categories[p.Category], p)
	}

	return utils.Response(c, http.StatusOK, "Categories with products fetched successfully", categories, nil, nil)
}

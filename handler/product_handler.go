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
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch products", nil, err, nil)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category); err != nil {
			return utils.Response(c, http.StatusInternalServerError, "Failed to parse products", nil, err, nil)
		}
		products = append(products, p)
	}
	return utils.Response(c, http.StatusOK, "Products fetched successfully", products, nil, nil)
}

func CreateProduct(c echo.Context) error {
	var p models.Product
	if err := c.Bind(&p); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, nil)
	}

	if err := validate.Struct(p); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, err, nil)
	}

	result, err := db.DB.Exec("INSERT INTO products (name, price, category) VALUES (?, ?, ?)", p.Name, p.Price, p.Category)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to create product", nil, err, nil)
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)
	return utils.Response(c, http.StatusCreated, "Product created successfully", p, nil, nil)
}

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	var p models.Product
	if err := c.Bind(&p); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, nil)
	}

	_, err := db.DB.Exec("UPDATE products SET name = ?, price = ?, category = ? WHERE id = ?", p.Name, p.Price, p.Category, id)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to update product", nil, err, nil)
	}

	return utils.Response(c, http.StatusOK, "Product updated successfully", p, nil, nil)
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")

	_, err := db.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to delete product", nil, err, nil)
	}

	return utils.Response(c, http.StatusOK, "Product deleted successfully", nil, nil, nil)
}

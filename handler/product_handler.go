package handler

import (
	"aro-shop/db"
	"aro-shop/models"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

func GetProducts(c echo.Context) error {
	rows, err := db.DB.Query("SELECT id, name, price, category FROM products")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Data:    nil,
			Message: "Failed to fetch products",
			Errors:  []string{err.Error()},
			ErrorID: "500-001",
		})
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category); err != nil {
			return c.JSON(http.StatusInternalServerError, models.Response{
				Data:    nil,
				Message: "Failed to parse products",
				Errors:  []string{err.Error()},
				ErrorID: "500-002",
			})
		}
		products = append(products, p)
	}
	return c.JSON(http.StatusOK, models.Response{
		Data:    products,
		Message: "Products fetched successfully",
		Errors:  nil,
	})
}

func CreateProduct(c echo.Context) error {
	var p models.Product
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Data:    nil,
			Message: "Invalid request",
			Errors:  []string{err.Error()},
			ErrorID: "400-001",
		})
	}

	if err := validate.Struct(p); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, models.Response{
			Data:    nil,
			Message: "Validation failed",
			Errors:  []string{validationErrors.Error()},
			ErrorID: "400-002",
		})
	}

	result, err := db.DB.Exec("INSERT INTO products (name, price, category) VALUES (?, ?, ?)", p.Name, p.Price, p.Category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Data:    nil,
			Message: "Failed to create product",
			Errors:  []string{err.Error()},
			ErrorID: "500-003",
		})
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)
	return c.JSON(http.StatusCreated, models.Response{
		Data:    p,
		Message: "Product created successfully",
		Errors:  nil,
	})
}

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Data:    nil,
			Message: "Invalid product ID",
			Errors:  []string{err.Error()},
			ErrorID: "400-003",
		})
	}

	var p models.Product
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Data:    nil,
			Message: "Invalid request",
			Errors:  []string{err.Error()},
			ErrorID: "400-004",
		})
	}

	if err := validate.Struct(p); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, models.Response{
			Data:    nil,
			Message: "Validation failed",
			Errors:  []string{validationErrors.Error()},
			ErrorID: "400-005",
		})
	}

	result, err := db.DB.Exec("UPDATE products SET name = ?, price = ?, category = ? WHERE id = ?", p.Name, p.Price, p.Category, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Data:    nil,
			Message: "Failed to update product",
			Errors:  []string{err.Error()},
			ErrorID: "500-004",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{
			Data:    nil,
			Message: "Product not found",
			Errors:  nil,
			ErrorID: "404-001",
		})
	}

	p.ID = productID
	return c.JSON(http.StatusOK, models.Response{
		Data:    p,
		Message: "Product updated successfully",
		Errors:  nil,
	})
}

func GetProductByID(c echo.Context) error {
	id := c.Param("id")

	var p models.Product
	err := db.DB.QueryRow("SELECT id, name, price, category FROM products WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Price, &p.Category)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Response{
			Data:    nil,
			Message: "Product not found",
			Errors:  []string{err.Error()},
			ErrorID: "404-003",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Data:    p,
		Message: "Product retrieved successfully",
		Errors:  nil,
	})
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")

	result, err := db.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Data:    nil,
			Message: "Failed to delete product",
			Errors:  []string{err.Error()},
			ErrorID: "500-005",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{
			Data:    nil,
			Message: "Product not found",
			Errors:  nil,
			ErrorID: "404-004",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Data:    nil,
		Message: "Product deleted successfully",
		Errors:  nil,
	})
}

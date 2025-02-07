package handler

import (
	"aro-shop/db"
	"aro-shop/models"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

func generateErrorID() *string {
	errorID := "ERR-" + time.Now().Format("150405")
	return &errorID
}

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
		return c.JSON(http.StatusInternalServerError, models.Response{
			Data:    nil,
			Message: "Failed to fetch products",
			Errors:  []map[string]string{{"error": err.Error()}},
			ErrorID: generateErrorID(),
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
				Errors:  []map[string]string{{"error": err.Error()}},
				ErrorID: generateErrorID(),
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
			Errors:  []map[string]string{{"error": err.Error()}},
			ErrorID: generateErrorID(),
		})
	}

	if err := validate.Struct(p); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, models.Response{
			Data:    nil,
			Message: "Validation failed",
			Errors:  []map[string]string{{"error": validationErrors.Error()}},
			ErrorID: generateErrorID(),
		})
	}

	result, err := db.DB.Exec("INSERT INTO products (name, price, category) VALUES (?, ?, ?)", p.Name, p.Price, p.Category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Data:    nil,
			Message: "Failed to create product",
			Errors:  []map[string]string{{"error": err.Error()}},
			ErrorID: generateErrorID(),
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

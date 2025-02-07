package handler

import (
	"aro-shop/db"
	"aro-shop/models"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func GetTransactions(c echo.Context) error {
	rows, err := db.DB.Query("SELECT id, product_id, quantity, total FROM transactions")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "Failed to fetch transactions", Errors: []string{err.Error()}})
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.ProductID, &t.Quantity, &t.Total); err != nil {
			return c.JSON(http.StatusInternalServerError, models.Response{Message: "Error scanning transactions", Errors: []string{err.Error()}})
		}
		transactions = append(transactions, t)
	}

	return c.JSON(http.StatusOK, models.Response{Data: transactions, Message: "Transactions retrieved successfully"})
}

func CreateTransaction(c echo.Context) error {
	var t models.Transaction
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "Invalid request", Errors: []string{err.Error()}})
	}

	if err := validate.Struct(t); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, models.Response{Message: "Validation failed", Errors: []string{validationErrors.Error()}})
	}

	var price float64
	err := db.DB.QueryRow("SELECT price FROM products WHERE id = ?", t.ProductID).Scan(&price)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Response{Message: "Product not found"})
	}

	t.Total = float64(t.Quantity) * price

	result, err := db.DB.Exec("INSERT INTO transactions (product_id, quantity, total) VALUES (?, ?, ?)", t.ProductID, t.Quantity, t.Total)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "Failed to create transaction", Errors: []string{err.Error()}})
	}

	transactionID, _ := result.LastInsertId()
	t.ID = int(transactionID)
	return c.JSON(http.StatusCreated, models.Response{Data: t, Message: "Transaction created successfully"})
}

func UpdateTransaction(c echo.Context) error {
	id := c.Param("id")

	transactionID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "Invalid transaction ID"})
	}

	var t models.Transaction
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "Invalid request", Errors: []string{err.Error()}})
	}

	if err := validate.Struct(t); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, models.Response{Message: "Validation failed", Errors: []string{validationErrors.Error()}})
	}

	var price float64
	err = db.DB.QueryRow("SELECT price FROM products WHERE id = ?", t.ProductID).Scan(&price)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Response{Message: "Product not found"})
	}

	t.Total = float64(t.Quantity) * price

	result, err := db.DB.Exec("UPDATE transactions SET product_id = ?, quantity = ?, total = ? WHERE id = ?", t.ProductID, t.Quantity, t.Total, transactionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "Failed to update transaction", Errors: []string{err.Error()}})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{Message: "Transaction not found"})
	}

	t.ID = transactionID
	return c.JSON(http.StatusOK, models.Response{Data: t, Message: "Transaction updated successfully"})
}

func DeleteTransaction(c echo.Context) error {
	id := c.Param("id")

	transactionID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "Invalid transaction ID"})
	}

	result, err := db.DB.Exec("DELETE FROM transactions WHERE id = ?", transactionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "Failed to delete transaction", Errors: []string{err.Error()}})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{Message: "Transaction not found"})
	}

	return c.JSON(http.StatusOK, models.Response{Message: "Transaction deleted successfully"})
}

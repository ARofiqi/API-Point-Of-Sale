package handler

import (
	"aro-shop/db"
	"aro-shop/models"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func GetTransactions(c echo.Context) error {
	rows, err := db.DB.Query("SELECT id, date, total FROM transactions")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to fetch transactions",
			Errors:  []string{err.Error()},
		})
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var dateBytes []byte

		if err := rows.Scan(&t.ID, &dateBytes, &t.Total); err != nil {
			return c.JSON(http.StatusInternalServerError, models.Response{
				Message: "Error scanning transactions",
				Errors:  []string{err.Error()},
			})
		}

		t.Date = string(dateBytes)

		items, err := getTransactionItems(t.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.Response{
				Message: "Error fetching transaction items",
				Errors:  []string{err.Error()},
			})
		}
		t.Items = items

		transactions = append(transactions, t)
	}

	return c.JSON(http.StatusOK, models.Response{
		Data:    transactions,
		Message: "Transactions retrieved successfully",
	})
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

	if len(t.Items) == 0 {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "Transaction must contain at least one item"})
	}

	t.Date = time.Now().Format("2006-01-02 15:04:05")

	tx, err := db.DB.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to start transaction",
			Errors:  []string{err.Error()},
		})
	}

	result, err := tx.Exec("INSERT INTO transactions (date, total) VALUES (?, ?)", t.Date, 0)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to create transaction",
			Errors:  []string{err.Error()},
		})
	}
	transactionID, _ := result.LastInsertId()

	var total float64
	for i, item := range t.Items {
		var price float64
		err := tx.QueryRow("SELECT price FROM products WHERE id = ?", item.ProductID).Scan(&price)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusNotFound, models.Response{Message: "Product not found", Errors: []string{err.Error()}})
		}
		subTotal := float64(item.Quantity) * price
		total += subTotal

		result, err := tx.Exec("INSERT INTO transaction_items (transaction_id, product_id, quantity, sub_total) VALUES (?, ?, ?, ?)", transactionID, item.ProductID, item.Quantity, subTotal)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, models.Response{Message: "Failed to create transaction items", Errors: []string{err.Error()}})
		}

		itemID, _ := result.LastInsertId()
		t.Items[i].ID = int(itemID)
		t.Items[i].TransactionID = int(transactionID)
		t.Items[i].SubTotal = subTotal
	}

	_, err = tx.Exec("UPDATE transactions SET total = ? WHERE id = ?", total, transactionID)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "Failed to update transaction total", Errors: []string{err.Error()}})
	}

	tx.Commit()
	t.ID = int(transactionID)
	t.Total = total
	return c.JSON(http.StatusCreated, models.Response{Data: t, Message: "Transaction created successfully"})
}

func getTransactionItems(transactionID int) ([]models.TransactionItem, error) {
	rows, err := db.DB.Query("SELECT id, transaction_id, product_id, quantity, sub_total FROM transaction_items WHERE transaction_id = ?", transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TransactionItem
	for rows.Next() {
		var item models.TransactionItem
		if err := rows.Scan(&item.ID, &item.TransactionID, &item.ProductID, &item.Quantity, &item.SubTotal); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

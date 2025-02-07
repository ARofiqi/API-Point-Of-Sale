package handler

import (
	"aro-shop/db"
	"aro-shop/models"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func generateErrorID() *string {
	errorID := "ERR-" + time.Now().Format("150405")
	return &errorID
}

func handleError(c echo.Context, statusCode int, message string, err error) error {
	errorID := time.Now().UnixNano()
	log.Printf("[ERROR %d] %s: %v", errorID, message, err)

	return c.JSON(statusCode, models.Response{
		Data:    nil,
		Message: message,
		Errors:  []map[string]string{{"error_code": "ERR-" + time.Now().Format("150405")}},
		ErrorID: generateErrorID(),
	})
}

func GetTransactions(c echo.Context) error {
	rows, err := db.DB.Query("SELECT id, date, total FROM transactions")
	if err != nil {
		return handleError(c, http.StatusInternalServerError, "Failed to fetch transactions", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var dateBytes []byte

		if err := rows.Scan(&t.ID, &dateBytes, &t.Total); err != nil {
			return handleError(c, http.StatusInternalServerError, "Error scanning transactions", err)
		}

		t.Date = string(dateBytes)

		items, err := getTransactionItems(t.ID)
		if err != nil {
			return handleError(c, http.StatusInternalServerError, "Error fetching transaction items", err)
		}
		t.Items = items

		transactions = append(transactions, t)
	}

	return c.JSON(http.StatusOK, models.Response{
		Data:    transactions,
		Message: "Transactions retrieved successfully",
		Errors:  nil,
		ErrorID: nil,
	})
}

func CreateTransaction(c echo.Context) error {
	var t models.Transaction
	if err := c.Bind(&t); err != nil {
		return handleError(c, http.StatusBadRequest, "Invalid request format", err)
	}

	if err := validate.Struct(t); err != nil {
		return handleError(c, http.StatusBadRequest, "Validation failed", err)
	}

	if len(t.Items) == 0 {
		return c.JSON(http.StatusBadRequest, models.Response{
			Data:    nil,
			Message: "Transaction must contain at least one item",
			Errors:  []map[string]string{{"items": "At least one item is required"}},
			ErrorID: nil,
		})
	}

	t.Date = time.Now().Format("2006-01-02 15:04:05")

	tx, err := db.DB.Begin()
	if err != nil {
		return handleError(c, http.StatusInternalServerError, "Failed to start transaction", err)
	}

	result, err := tx.Exec("INSERT INTO transactions (date, total) VALUES (?, ?)", t.Date, 0)
	if err != nil {
		tx.Rollback()
		return handleError(c, http.StatusInternalServerError, "Failed to create transaction", err)
	}
	transactionID, _ := result.LastInsertId()

	var total float64
	for i, item := range t.Items {
		var price float64
		err := tx.QueryRow("SELECT price FROM products WHERE id = ?", item.ProductID).Scan(&price)
		if err != nil {
			tx.Rollback()
			return handleError(c, http.StatusNotFound, "Product not found", err)
		}
		subTotal := float64(item.Quantity) * price
		total += subTotal

		result, err := tx.Exec("INSERT INTO transaction_items (transaction_id, product_id, quantity, sub_total) VALUES (?, ?, ?, ?)", transactionID, item.ProductID, item.Quantity, subTotal)
		if err != nil {
			tx.Rollback()
			return handleError(c, http.StatusInternalServerError, "Failed to create transaction items", err)
		}

		itemID, _ := result.LastInsertId()
		t.Items[i].ID = int(itemID)
		t.Items[i].TransactionID = int(transactionID)
		t.Items[i].SubTotal = subTotal
	}

	_, err = tx.Exec("UPDATE transactions SET total = ? WHERE id = ?", total, transactionID)
	if err != nil {
		tx.Rollback()
		return handleError(c, http.StatusInternalServerError, "Failed to update transaction total", err)
	}

	tx.Commit()
	t.ID = int(transactionID)
	t.Total = total
	return c.JSON(http.StatusCreated, models.Response{
		Data:    t,
		Message: "Transaction created successfully",
		Errors:  nil,
		ErrorID: nil,
	})
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

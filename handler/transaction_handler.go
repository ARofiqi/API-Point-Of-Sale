package handler

import (
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetTransactions(c echo.Context) error {
	rows, err := db.DB.Query("SELECT id, date, total FROM transactions")
	if err != nil {
		errorDetails := map[string]string{"query_error": "Gagal mengambil transaksi"}
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch transactions", nil, err, errorDetails)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var dateBytes []byte

		if err := rows.Scan(&t.ID, &dateBytes, &t.Total); err != nil {
			errorDetails := map[string]string{"scan_error": "Gagal membaca data transaksi"}
			return utils.Response(c, http.StatusInternalServerError, "Error scanning transactions", nil, err, errorDetails)
		}

		t.Date = string(dateBytes)

		items, err := getTransactionItems(t.ID)
		if err != nil {
			errorDetails := map[string]string{"fetch_items_error": "Gagal mengambil item transaksi"}
			return utils.Response(c, http.StatusInternalServerError, "Error fetching transaction items", nil, err, errorDetails)
		}
		t.Items = items

		transactions = append(transactions, t)
	}

	return utils.Response(c, http.StatusOK, "Transactions retrieved successfully", transactions, nil, nil)
}

func CreateTransaction(c echo.Context) error {
	var t models.Transaction
	if err := c.Bind(&t); err != nil {
		errorDetails := map[string]string{"binding_error": "Format permintaan tidak valid"}
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, errorDetails)
	}

	if err := validate.Struct(t); err != nil {
		errorDetails := map[string]string{"validation_error": "Validasi gagal"}
		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, err, errorDetails)
	}

	if len(t.Items) == 0 {
		errorDetails := map[string]string{"items": "At least one item is required"}
		return utils.Response(c, http.StatusBadRequest, "Transaction must contain at least one item", nil, nil, errorDetails)
	}

	t.Date = time.Now().Format("2006-01-02 15:04:05")

	tx, err := db.DB.Begin()
	if err != nil {
		errorDetails := map[string]string{"transaction_error": "Gagal memulai transaksi"}
		return utils.Response(c, http.StatusInternalServerError, "Failed to start transaction", nil, err, errorDetails)
	}

	result, err := tx.Exec("INSERT INTO transactions (date, total) VALUES (?, ?)", t.Date, 0)
	if err != nil {
		tx.Rollback()
		errorDetails := map[string]string{"insert_transaction_error": "Gagal membuat transaksi"}
		return utils.Response(c, http.StatusInternalServerError, "Failed to create transaction", nil, err, errorDetails)
	}
	transactionID, _ := result.LastInsertId()

	var total float64
	for i, item := range t.Items {
		var price float64
		err := tx.QueryRow("SELECT price FROM products WHERE id = ?", item.ProductID).Scan(&price)
		if err != nil {
			tx.Rollback()
			errorDetails := map[string]string{"product_not_found": fmt.Sprintf("Produk dengan ID %d tidak ditemukan", item.ProductID)}
			return utils.Response(c, http.StatusNotFound, "Product not found", nil, err, errorDetails)
		}
		subTotal := float64(item.Quantity) * price
		total += subTotal

		result, err := tx.Exec("INSERT INTO transaction_items (transaction_id, product_id, quantity, sub_total) VALUES (?, ?, ?, ?)", transactionID, item.ProductID, item.Quantity, subTotal)
		if err != nil {
			tx.Rollback()
			errorDetails := map[string]string{"insert_item_error": "Gagal menambahkan item transaksi"}
			return utils.Response(c, http.StatusInternalServerError, "Failed to create transaction items", nil, err, errorDetails)
		}

		itemID, _ := result.LastInsertId()
		t.Items[i].ID = int(itemID)
		t.Items[i].TransactionID = int(transactionID)
		t.Items[i].SubTotal = subTotal
	}

	_, err = tx.Exec("UPDATE transactions SET total = ? WHERE id = ?", total, transactionID)
	if err != nil {
		tx.Rollback()
		errorDetails := map[string]string{"update_total_error": "Gagal memperbarui total transaksi"}
		return utils.Response(c, http.StatusInternalServerError, "Failed to update transaction total", nil, err, errorDetails)
	}

	tx.Commit()
	t.ID = int(transactionID)
	t.Total = total

	return utils.Response(c, http.StatusCreated, "Transaction created successfully", t, nil, nil)
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

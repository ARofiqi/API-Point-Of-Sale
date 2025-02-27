package handler

import (
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetTransactions(c echo.Context) error {
	var transactions []models.Transaction
	if err := db.DB.Preload("Items").Find(&transactions).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch transactions", nil, err, nil)
	}
	return utils.Response(c, http.StatusOK, "Transactions retrieved successfully", transactions, nil, nil)
}

func CreateTransaction(c echo.Context) error {
	var (
		t            models.Transaction
		errorDetails = make(map[string]string)
	)

	if err := c.Bind(&t); err != nil {
		errorDetails = utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Invalid request format", nil, err, errorDetails)
	}

	if err := validate.Struct(t); err != nil {
		errorDetails = utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, err, errorDetails)
	}

	if len(t.Items) == 0 {
		errorDetails["items"] = "Transaction must contain at least one item"
		return utils.Response(c, http.StatusBadRequest, "Validation failed", nil, nil, errorDetails)
	}

	t.Date = time.Now()
	tx := db.DB.Begin()
	if err := tx.Create(&t).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, http.StatusInternalServerError, "Failed to create transaction", nil, err, nil)
	}

	var total float64
	for i := range t.Items {
		var price float64
		if err := tx.Model(&models.Product{}).Select("price").Where("id = ?", t.Items[i].ProductID).Scan(&price).Error; err != nil {
			tx.Rollback()
			return utils.Response(c, http.StatusNotFound, "Product not found", nil, err, nil)
		}

		t.Items[i].SubTotal = float64(t.Items[i].Quantity) * price
		t.Items[i].TransactionID = t.ID
		total += t.Items[i].SubTotal
	}

	if err := tx.Create(&t.Items).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, http.StatusInternalServerError, "Failed to create transaction items", nil, err, nil)
	}

	t.Total = total
	if err := tx.Save(&t).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, http.StatusInternalServerError, "Failed to update transaction total", nil, err, nil)
	}
	tx.Commit()

	return utils.Response(c, http.StatusCreated, "Transaction created successfully", t, nil, nil)
}

func GetTransactionSubtotal(c echo.Context) error {
	var (
		transaction  models.Transaction
		errorDetails = make(map[string]string)
	)

	transactionID := c.Param("id")
	if err := db.DB.Preload("Items").First(&transaction, transactionID).Error; err != nil {
		errorDetails = utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusNotFound, "Transaction not found", nil, err, errorDetails)
	}

	subtotal := 0.0
	for _, item := range transaction.Items {
		subtotal += item.SubTotal
	}

	result := map[string]interface{}{
		"transaction_id": transactionID,
		"subtotal":       subtotal,
	}

	return utils.Response(c, http.StatusOK, "Transaction subtotal retrieved successfully", result, nil, nil)
}

func GetTransactionsByDateRange(c echo.Context) error {

	var (
		errorDetails = make(map[string]string)
		startDate    = c.QueryParam("start")
		endDate      = c.QueryParam("end")
	)

	if startDate == "" || endDate == "" {
		errorDetails["date_range_error"] = "Start date dan end date diperlukan"
		return utils.Response(c, http.StatusBadRequest, "Start date and end date are required", nil, nil, errorDetails)
	}

	var transactions []models.Transaction
	if err := db.DB.Preload("Items").Where("date BETWEEN ? AND ?", startDate, endDate).Find(&transactions).Error; err != nil {
		// errorDetails["query_error"] = "Gagal mengambil transaksi berdasarkan rentang tanggal"
		errorDetails = utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch transactions by date range", nil, err, errorDetails)
	}

	return utils.Response(c, http.StatusOK, "Transactions retrieved successfully", transactions, nil, nil)
}

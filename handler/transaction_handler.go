package handler

import (
	"aro-shop/cache"
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/queue"
	"aro-shop/utils"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetTransactions(c echo.Context) error {
	cacheKey := "all_transactions"

	cachedData, err := cache.GetCache(cacheKey)
	if err == nil {
		var transactions []models.Transaction
		if json.Unmarshal([]byte(cachedData), &transactions) == nil {
			return utils.Response(c, http.StatusOK, "Transactions retrieved successfully (from cache)", transactions, nil, nil)
		}
	}

	var transactions []models.Transaction
	if err := db.DB.Preload("Items").Find(&transactions).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch transactions", nil, err, nil)
	}

	dataJSON, _ := json.Marshal(transactions)
	cache.SetCache(cacheKey, string(dataJSON), 5*time.Minute)

	return utils.Response(c, http.StatusOK, "Transactions retrieved successfully", transactions, nil, nil)
}

func CreateTransaction(c echo.Context) error {
	var (
		t            models.Transaction
		errorDetails = make(map[string]string)
	)

	if err := c.Bind(&t); err != nil {
		errorDetails = utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Format permintaan tidak valid", nil, err, errorDetails)
	}

	if err := validate.Struct(t); err != nil {
		errorDetails = utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Validasi gagal", nil, err, errorDetails)
	}

	if len(t.Items) == 0 {
		errorDetails["items"] = "Transaksi harus memiliki setidaknya satu item"
		return utils.Response(c, http.StatusBadRequest, "Validasi gagal", nil, nil, errorDetails)
	}

	t.Date = time.Now()

	transactionJSON, err := json.Marshal(t)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Gagal serialisasi transaksi", nil, err, nil)
	}

	if err := queue.PublishTransaction(transactionJSON); err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Gagal mengirim transaksi ke antrian", nil, err, nil)
	}

	notificationMessage := "Transaksi baru telah dibuat"
	if err := queue.PublishNotification(notificationMessage); err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Gagal mengirim notifikasi", nil, err, nil)
	}

	cache.DeleteCache("all_transactions")

	return utils.Response(c, http.StatusAccepted, "Transaksi berhasil dikirim ke antrian", nil, nil, nil)
}

func GetTransactionSubtotal(c echo.Context) error {
	transactionID := c.Param("id")
	cacheKey := "transaction_subtotal_" + transactionID

	cachedData, err := cache.GetCache(cacheKey)
	if err == nil {
		var result map[string]interface{}
		if json.Unmarshal([]byte(cachedData), &result) == nil {
			return utils.Response(c, http.StatusOK, "Transaction subtotal retrieved successfully (from cache)", result, nil, nil)
		}
	}

	var transaction models.Transaction
	if err := db.DB.Preload("Items").First(&transaction, transactionID).Error; err != nil {
		return utils.Response(c, http.StatusNotFound, "Transaction not found", nil, err, nil)
	}

	subtotal := 0.0
	for _, item := range transaction.Items {
		subtotal += item.SubTotal
	}

	result := map[string]interface{}{
		"transaction_id": transactionID,
		"subtotal":       subtotal,
	}

	dataJSON, _ := json.Marshal(result)
	cache.SetCache(cacheKey, string(dataJSON), 5*time.Minute)

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

	cacheKey := "transactions_" + startDate + "_" + endDate

	cachedData, err := cache.GetCache(cacheKey)
	if err == nil {
		var transactions []models.Transaction
		if json.Unmarshal([]byte(cachedData), &transactions) == nil {
			return utils.Response(c, http.StatusOK, "Transactions retrieved successfully (from cache)", transactions, nil, nil)
		}
	}

	var transactions []models.Transaction
	if err := db.DB.Preload("Items").Where("date BETWEEN ? AND ?", startDate, endDate).Find(&transactions).Error; err != nil {
		errorDetails = utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch transactions by date range", nil, err, errorDetails)
	}

	dataJSON, _ := json.Marshal(transactions)
	cache.SetCache(cacheKey, string(dataJSON), 5*time.Minute)

	return utils.Response(c, http.StatusOK, "Transactions retrieved successfully", transactions, nil, nil)
}


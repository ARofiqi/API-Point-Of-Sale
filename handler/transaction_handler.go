package handler

import (
	"aro-shop/cache"
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/queue"
	"aro-shop/utils"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	cachedDataTransactions = []string{
		"all_transactions",
	}
)

func GetTransactions(c echo.Context) error {
	cacheKeyPrefix := "transactions_page_"
	errorDetails := make(models.ErrorDetails)

	// Ambil parameter page dan limit dari query string
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	cacheKey := fmt.Sprintf("%s%d_%d", cacheKeyPrefix, page, limit)

	// Cek apakah data ada di Redis
	cachedData, err := cache.GetCache(cacheKey)
	if err == nil {
		var transactions []models.Transaction
		if json.Unmarshal([]byte(cachedData), &transactions) == nil {
			return utils.Response(c, http.StatusOK, "Transactions retrieved successfully (from cache)", transactions, nil, nil)
		}
	}

	// Ambil total data untuk indexing
	var total int64
	if err := db.DB.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		errorDetails["database"] = "Failed to count transactions"
		return utils.Response(c, http.StatusInternalServerError, "Internal server error", nil, err, errorDetails)
	}

	// Jika tidak ada di Redis, ambil dari database dengan pagination
	var transactions []models.Transaction
	if err := db.DB.Preload("Items").Limit(limit).Offset(offset).Find(&transactions).Error; err != nil {
		errorDetails["database"] = "Failed to fetch transactions"
		return utils.Response(c, http.StatusInternalServerError, "Internal server error", nil, err, errorDetails)
	}

	// Simpan hasil query ke Redis
	dataJSON, _ := json.Marshal(transactions)
	cache.SetCache(cacheKey, string(dataJSON), 10*time.Minute)

	// Struktur respons dengan pagination
	response := map[string]interface{}{
		"transactions": transactions,
		"pagination": map[string]interface{}{
			"current_page": page,
			"per_page":     limit,
			"total_data":   total,
			"total_pages":  int(math.Ceil(float64(total) / float64(limit))),
		},
	}

	return utils.Response(c, http.StatusOK, "Transactions retrieved successfully", response, nil, nil)
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
	cache.SetCache(cacheKey, string(dataJSON), 10*time.Minute)

	return utils.Response(c, http.StatusOK, "Transaction subtotal retrieved successfully", result, nil, nil)
}

func GetTransactionsByDateRange(c echo.Context) error {
	var (
		errorDetails = make(map[string]string)
		startDate    = c.QueryParam("start")
		endDate      = c.QueryParam("end")
		page, _      = strconv.Atoi(c.QueryParam("page"))
		limit, _     = strconv.Atoi(c.QueryParam("limit"))
	)

	if startDate == "" || endDate == "" {
		errorDetails["date_range_error"] = "Start date dan end date diperlukan"
		return utils.Response(c, http.StatusBadRequest, "Start date and end date are required", nil, nil, errorDetails)
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	cacheKey := fmt.Sprintf("transactions_%s_%s_page_%d_limit_%d", startDate, endDate, page, limit)

	cachedData, err := cache.GetCache(cacheKey)
	if err == nil {
		var cachedResponse struct {
			Pagination   map[string]interface{} `json:"pagination"`
			Transactions []models.Transaction   `json:"transactions"`
		}
		if json.Unmarshal([]byte(cachedData), &cachedResponse) == nil {
			return utils.Response(c, http.StatusOK, "Transactions retrieved successfully (from cache)", cachedResponse, nil, nil)
		}
	}

	var transactions []models.Transaction
	var totalRecords int64

	db.DB.Model(&models.Transaction{}).Where("date BETWEEN ? AND ?", startDate, endDate).Count(&totalRecords)

	if err := db.DB.Preload("Items").Where("date BETWEEN ? AND ?", startDate, endDate).Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		errorDetails = utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusInternalServerError, "Failed to fetch transactions by date range", nil, err, errorDetails)
	}

	pagination := map[string]interface{}{
		"current_page": page,
		"per_page":     limit,
		"total_data":   totalRecords,
		"total_pages":  int(math.Ceil(float64(totalRecords) / float64(limit))),
	}

	responseData := map[string]interface{}{
		"pagination":   pagination,
		"transactions": transactions,
	}

	dataJSON, _ := json.Marshal(responseData)
	cache.SetCache(cacheKey, string(dataJSON), 10*time.Minute)

	return utils.Response(c, http.StatusOK, "Transactions retrieved successfully", responseData, nil, nil)
}

func CreateTransaction(c echo.Context) error {
	var (
		t            models.Transaction
		errorDetails = make(map[string]string)
	)

	if err := c.Bind(&t); err != nil {
		// errorDetails = utils.ParseValidationErrors(err)
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

	go cache.ResetRedisCache(cachedDataTransactions...)

	return utils.Response(c, http.StatusAccepted, "Transaksi berhasil dikirim ke antrian", nil, nil, nil)
}

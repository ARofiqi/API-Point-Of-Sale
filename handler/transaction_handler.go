package handler

import (
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

	// Bind request body ke struct
	if err := c.Bind(&t); err != nil {
		errorDetails = utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Format permintaan tidak valid", nil, err, errorDetails)
	}

	// Validasi input
	if err := validate.Struct(t); err != nil {
		errorDetails = utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Validasi gagal", nil, err, errorDetails)
	}

	if len(t.Items) == 0 {
		errorDetails["items"] = "Transaksi harus memiliki setidaknya satu item"
		return utils.Response(c, http.StatusBadRequest, "Validasi gagal", nil, nil, errorDetails)
	}

	// Set tanggal transaksi
	t.Date = time.Now()

	// Serialize transaksi ke JSON
	transactionJSON, err := json.Marshal(t)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Gagal serialisasi transaksi", nil, err, nil)
	}

	// Kirim transaksi ke RabbitMQ
	if err := queue.PublishTransaction(transactionJSON); err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Gagal mengirim transaksi ke antrian", nil, err, nil)
	}

	// Kirim notifikasi ke RabbitMQ
	notificationMessage := "Transaksi baru telah dibuat"
	if err := queue.PublishNotification(notificationMessage); err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Gagal mengirim notifikasi", nil, err, nil)
	}

	return utils.Response(c, http.StatusAccepted, "Transaksi berhasil dikirim ke antrian", nil, nil, nil)
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

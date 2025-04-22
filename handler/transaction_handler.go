package handler

import (
	"aro-shop/cache"
	"aro-shop/db"
	"aro-shop/dto"
	"aro-shop/models"
	"aro-shop/queue"
	"aro-shop/utils"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	cachedDataTransactions = []string{
		"all_transactions",
	}
)

func GetTransactions(c echo.Context) error {
	cacheKeyPrefix := "transactions_page_"
	errorDetails := make(dto.ErrorDetails)

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
	if err := db.DB.
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Payment").
		Preload("Payment.PaymentMethod").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		errorDetails["database"] = "Failed to fetch transactions"
		return utils.Response(c, http.StatusInternalServerError, "Internal server error", nil, err, errorDetails)
	}

	// Map to DTO response
	transactionResponses := MapTransactionsToResponse(transactions)

	// Simpan hasil query ke Redis
	dataJSON, _ := json.Marshal(transactionResponses)
	cache.SetCache(cacheKey, string(dataJSON), 10*time.Minute)

	// Struktur respons dengan pagination
	response := map[string]interface{}{
		"transactions": transactionResponses,
		"pagination": map[string]interface{}{
			"current_page": page,
			"per_page":     limit,
			"total_data":   total,
			"total_pages":  int(math.Ceil(float64(total) / float64(limit))),
		},
	}

	return utils.Response(c, http.StatusOK, "Transactions retrieved successfully", response, nil, nil)
}

func GetTransactionsById(c echo.Context) error {
	cacheKeyPrefix := "transactions_"
	errorDetails := make(dto.ErrorDetails)

	// Ambil parameter id dari URL params
	id := c.Param("id")

	// Validasi format UUID (opsional tapi direkomendasikan)
	_, err := uuid.Parse(id)
	if err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid UUID format", nil, nil, nil)
	}

	cacheKey := fmt.Sprintf("%s%s", cacheKeyPrefix, id)

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
	var transaction models.Transaction
	if err := db.DB.
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Payment").
		Preload("Payment.PaymentMethod").
		Find(&transaction).Error; err != nil {
		errorDetails["database"] = "Failed to fetch transactions"
		return utils.Response(c, http.StatusInternalServerError, "Internal server error", nil, err, errorDetails)
	}

	TransactionResponse := dto.TransactionResponse{
		ID:         transaction.ID,
		User:       dto.SimpleUserResponse{ID: transaction.User.ID, Name: transaction.User.Name, Email: transaction.User.Email},
		Date:       transaction.Date,
		AmountPaid: transaction.AmountPaid,
		Items:      MapTransactionItemToResponse(transaction.Items),
		Payment: &dto.PaymentResponse{
			ID:            transaction.Payment.ID,
			Status:        transaction.Payment.PaymentStatus,
			PaymentMethod: &dto.PaymentMethodSimple{ID: transaction.Payment.PaymentMethod.ID, Name: transaction.Payment.PaymentMethod.Name},
			PaidAt:        transaction.Payment.PaidAt,
			AmountPaid:    transaction.Payment.AmountPaid,
			CreatedAt:     transaction.Payment.CreatedAt,
			UpdatedAt:     transaction.Payment.UpdatedAt,
		},
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
	}

	// Simpan hasil query ke Redis
	dataJSON, _ := json.Marshal(TransactionResponse)
	cache.SetCache(cacheKey, string(dataJSON), 10*time.Minute)

	return utils.Response(c, http.StatusOK, "Transactions retrieved successfully", TransactionResponse, nil, nil)
}

func GetTransactionSubtotal(c echo.Context) error {
	transactionID := c.Param("id")
	cacheKey := "transaction_subtotal_" + transactionID

	// cachedData, err := cache.GetCache(cacheKey)
	// if err == nil {
	// 	var result map[string]interface{}
	// 	if json.Unmarshal([]byte(cachedData), &result) == nil {
	// 		return utils.Response(c, http.StatusOK, "Transaction subtotal retrieved successfully (from cache)", result, nil, nil)
	// 	}
	// }

	// Validasi UUID terlebih dahulu
	if _, err := uuid.Parse(transactionID); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid transaction ID format", nil, err, nil)
	}

	var transaction models.Transaction
	if err := db.DB.Preload("Items").Where("id = ?", transactionID).First(&transaction).Error; err != nil {
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
		req          dto.TransactionRequest
		errorDetails = make(map[string]string)
	)

	// Bind request ke struct DTO
	if err := c.Bind(&req); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Format permintaan tidak valid", nil, err, nil)
	}

	// Validasi menggunakan validator
	if err := validate.Struct(req); err != nil {
		errorDetails = utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Validasi gagal", nil, err, errorDetails)
	}

	// Validasi item produk
	if len(req.Items) == 0 {
		errorDetails["items"] = "Transaksi harus memiliki setidaknya satu item"
		return utils.Response(c, http.StatusBadRequest, "Validasi gagal", nil, nil, errorDetails)
	}

	var total float64
	var transactionItems []models.TransactionItem

	for _, item := range req.Items {
		var product models.Product
		if err := db.DB.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				errorDetails["product_id"] = fmt.Sprintf("Produk dengan ID %v tidak ditemukan", item.ProductID)
				return utils.Response(c, http.StatusBadRequest, "Produk tidak valid", nil, nil, errorDetails)
			}
			return utils.Response(c, http.StatusInternalServerError, "Gagal memeriksa produk", nil, err, nil)
		}

		subTotal := product.Price * float64(item.Quantity)
		total += subTotal

		transactionItems = append(transactionItems, models.TransactionItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			SubTotal:  subTotal,
		})
	}

	userID := c.Get("user_id")

	strID, ok := userID.(string)
	if !ok {
		return utils.Response(c, http.StatusUnauthorized, "User ID tidak valid", nil, nil, nil)
	}

	uid, err := uuid.Parse(strID)
	if err != nil || uid == uuid.Nil {
		return utils.Response(c, http.StatusUnauthorized, "User UUID tidak valid", nil, nil, nil)
	}

	transaction := models.Transaction{
		UserID:     uid,
		Date:       time.Now(),
		AmountPaid: total,
	}

	// Simpan transaksi terlebih dahulu untuk mendapatkan ID
	if err := db.DB.Create(&transaction).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Gagal menyimpan transaksi", nil, err, nil)
	}

	// Set TransactionID untuk setiap item
	for i := range transactionItems {
		transactionItems[i].TransactionID = transaction.ID
	}

	// Simpan semua item ke database
	if err := db.DB.Create(&transactionItems).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Gagal menyimpan item transaksi", nil, err, nil)
	}

	// Buat entitas pembayaran
	payment := models.Payment{
		TransactionID:   transaction.ID,
		PaymentMethodID: req.PaymentMethodID,
		PaymentStatus:   "pending",
		PaidAt:          nil,
	}

	if err := db.DB.Create(&payment).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Gagal menyimpan pembayaran", nil, err, nil)
	}

	if err := db.DB.
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Payment").
		Preload("Payment.PaymentMethod").
		First(&transaction, "id = ?", transaction.ID).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Gagal mengambil data transaksi lengkap", nil, err, nil)
	}

	// Kirim ke antrian
	// transactionJSON, err := json.Marshal(transaction)
	// if err != nil {
	// 	return utils.Response(c, http.StatusInternalServerError, "Gagal serialisasi transaksi", nil, err, nil)
	// }
	// if err := queue.PublishTransaction(transactionJSON); err != nil {
	// 	return utils.Response(c, http.StatusInternalServerError, "Gagal mengirim transaksi ke antrian", nil, err, nil)
	// }

	// Kirim notifikasi
	if err := queue.PublishNotification("Transaksi baru telah dibuat"); err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Gagal mengirim notifikasi", nil, err, nil)
	}

	// Reset Redis cache
	go cache.ResetRedisCache(cachedDataTransactions...)

	TransactionResponse := dto.TransactionResponse{
		ID:         transaction.ID,
		User:       dto.SimpleUserResponse{ID: transaction.User.ID, Name: transaction.User.Name, Email: transaction.User.Email},
		Date:       transaction.Date,
		AmountPaid: transaction.AmountPaid,
		Items:      MapTransactionItemToResponse(transaction.Items),
		Payment: &dto.PaymentResponse{
			ID:            transaction.Payment.ID,
			Status:        transaction.Payment.PaymentStatus,
			PaymentMethod: &dto.PaymentMethodSimple{ID: transaction.Payment.PaymentMethod.ID, Name: transaction.Payment.PaymentMethod.Name},
			PaidAt:        transaction.Payment.PaidAt,
			AmountPaid:    transaction.Payment.AmountPaid,
			CreatedAt:     transaction.Payment.CreatedAt,
			UpdatedAt:     transaction.Payment.UpdatedAt,
		},
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
	}

	return utils.Response(c, http.StatusCreated, "Transaksi berhasil dibuat", TransactionResponse, nil, nil)
}

func UpdateTransaction(c echo.Context) error {
	transactionID := c.Param("id")
	var transaction models.Transaction

	// Ambil data transaksi dari DB beserta relasinya
	if err := db.DB.Preload("Items").Preload("Payment").First(&transaction, "id = ?", transactionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.Response(c, http.StatusNotFound, "Transaksi tidak ditemukan", nil, nil, nil)
		}
		return utils.Response(c, http.StatusInternalServerError, "Gagal mengambil transaksi", nil, err, nil)
	}

	// Cek jika sudah dibayar
	if transaction.Payment.PaymentStatus == "paid" {
		return utils.Response(c, http.StatusBadRequest, "Transaksi sudah dibayar", nil, nil, nil)
	}

	now := time.Now()

	// Update status pembayaran
	if err := db.DB.Model(&transaction.Payment).Updates(models.Payment{
		PaymentStatus: "paid",
		PaidAt:        &now,
		AmountPaid:    transaction.AmountPaid,
	}).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Gagal memperbarui pembayaran", nil, err, nil)
	}

	// Kurangi stok produk
	for _, item := range transaction.Items {
		if err := db.DB.Model(&models.Product{}).
			Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
			Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
			return utils.Response(c, http.StatusInternalServerError, "Gagal mengurangi stok produk", nil, err, nil)
		}
	}

	// Ambil ulang data lengkap
	if err := db.DB.
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Payment").
		Preload("Payment.PaymentMethod").
		First(&transaction, "id = ?", transaction.ID).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Gagal mengambil data transaksi lengkap", nil, err, nil)
	}

	// Response
	TransactionResponse := dto.TransactionResponse{
		ID:         transaction.ID,
		Payment: &dto.PaymentResponse{
			ID:            transaction.Payment.ID,
			Status:        transaction.Payment.PaymentStatus,
			PaymentMethod: &dto.PaymentMethodSimple{ID: transaction.Payment.PaymentMethod.ID, Name: transaction.Payment.PaymentMethod.Name},
			PaidAt:        transaction.Payment.PaidAt,
			AmountPaid:    transaction.Payment.AmountPaid,
			CreatedAt:     transaction.Payment.CreatedAt,
			UpdatedAt:     transaction.Payment.UpdatedAt,
		},
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
	}

	return utils.Response(c, http.StatusOK, "Transaksi berhasil diperbarui dan dibayar", TransactionResponse, nil, nil)
}

func MapTransactionItemToResponse(items []models.TransactionItem) []dto.TransactionItemResponse {
	var response []dto.TransactionItemResponse
	for _, item := range items {
		response = append(response, dto.TransactionItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			SubTotal:    item.SubTotal,
		})
	}
	return response
}

func MapTransactionsToResponse(transactions []models.Transaction) []dto.TransactionResponse {
	var responses []dto.TransactionResponse
	for _, transaction := range transactions {
		res := dto.TransactionResponse{
			ID:         transaction.ID,
			User:       dto.SimpleUserResponse{ID: transaction.User.ID, Name: transaction.User.Name, Email: transaction.User.Email},
			Date:       transaction.Date,
			AmountPaid: transaction.AmountPaid,
			// Items:      MapTransactionItemToResponse(transaction.Items),
			Payment: &dto.PaymentResponse{
				ID:            transaction.Payment.ID,
				Status:        transaction.Payment.PaymentStatus,
				PaymentMethod: &dto.PaymentMethodSimple{
					ID: transaction.Payment.PaymentMethod.ID, 
					Name: transaction.Payment.PaymentMethod.Name,
				},
				PaidAt:        transaction.Payment.PaidAt,
				AmountPaid:    transaction.Payment.AmountPaid,
				CreatedAt:     transaction.Payment.CreatedAt,
				UpdatedAt:     transaction.Payment.UpdatedAt,
			},
			CreatedAt: transaction.CreatedAt,
			UpdatedAt: transaction.UpdatedAt,
		}
		responses = append(responses, res)
	}
	return responses
}

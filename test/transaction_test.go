package test

import (
	"aro-shop/handler"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetTransactions(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	SetupMockDB()

	if assert.NoError(t, handler.GetTransactions(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCreateTransaction(t *testing.T) {
	e := echo.New()
	SetupMockDB()

	requestBody := map[string]interface{}{
		"items": []map[string]interface{}{
			{"product_id": 1, "quantity": 2},
		},
	}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.CreateTransaction(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestGetTransactionSubtotal(t *testing.T) {
	e := echo.New()
	SetupMockDB()

	req := httptest.NewRequest(http.MethodGet, "/transactions/1/subtotal", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, handler.GetTransactionSubtotal(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetTransactionsByDateRange(t *testing.T) {
	e := echo.New()
	SetupMockDB()

	startDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	req := httptest.NewRequest(http.MethodGet, "/transactions?start="+startDate+"&end="+endDate, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.GetTransactionsByDateRange(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

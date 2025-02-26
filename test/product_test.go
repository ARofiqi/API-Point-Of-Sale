package test

import (
	"aro-shop/db"
	"aro-shop/handler"
	"aro-shop/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var token = "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDA1NTY5MjgsInJvbGUiOiJ1c2VyIiwidXNlcl9pZCI6IjY5MmIxODEyLTYzMmEtNDRlZi05MWM5LTgzYmE2ZmY4MDgwNiJ9.JsA9xK9L-K99juqysXlWSpglndRuaafKrOHI_G2embU"

var mock sqlmock.Sqlmock

func SetupMockDB() {
	sqlDB, mockDB, _ := sqlmock.New()
	mock = mockDB

	gormDB, _ := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	db.DB = gormDB
}

func TestGetProducts(t *testing.T) {
	SetupMockDB()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.GetProducts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCreateProduct(t *testing.T) {
	SetupMockDB()
	e := echo.New()

	product := models.Product{
		Name:       "Test Product",
		Price:      10000,
		CategoryID: 1,
	}
	jsonBody, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.CreateProduct(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestUpdateProduct(t *testing.T) {
	SetupMockDB()
	e := echo.New()

	product := models.Product{
		Name:       "Updated Product",
		Price:      20000,
		CategoryID: 1,
	}
	jsonBody, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, handler.UpdateProduct(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteProduct(t *testing.T) {
	SetupMockDB()
	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, handler.DeleteProduct(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetProductByID(t *testing.T) {
	SetupMockDB()
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, handler.GetProductByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetCategoriesWithProducts(t *testing.T) {
	SetupMockDB()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	req.Header.Set("Authorization", token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.GetCategoriesWithProducts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

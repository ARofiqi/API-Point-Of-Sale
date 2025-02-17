// package test

// import (
// 	"aro-shop/handler"
// 	"aro-shop/models"
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/labstack/echo/v4"
// 	"github.com/stretchr/testify/assert"
// )

// var token = "Barier " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk0MzAzMjcsInJvbGUiOiJhZG1pbiIsInVzZXJfaWQiOiJhYmQ5MDI5ZC1lODE5LTExZWYtYmE1NC1kMGM1ZDMxODBiY2UifQ.HnemBD3tl5YK0xtli2ewxiiZrm2S-7MwgUu3EGWwkIk"

// func TestGetProducts(t *testing.T) {
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/products", nil)
// 	req.Header.Set("Authorization", token)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	if assert.NoError(t, handler.GetProducts(c)) {
// 		assert.Equal(t, http.StatusOK, rec.Code)
// 	}
// }

// func TestCreateProduct(t *testing.T) {
// 	e := echo.New()

// 	product := models.Product{
// 		Name:       "Test Product",
// 		Price:      10000,
// 		CategoryID: 1,
// 	}
// 	jsonBody, _ := json.Marshal(product)
// 	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/products", bytes.NewBuffer(jsonBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", token) // Tambahkan token di sini
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	if assert.NoError(t, handler.CreateProduct(c)) {
// 		assert.Equal(t, http.StatusCreated, rec.Code)
// 	}
// }

// func TestUpdateProduct(t *testing.T) {
// 	e := echo.New()

// 	product := models.Product{
// 		Name:       "Updated Product",
// 		Price:      20000,
// 		CategoryID: 1,
// 	}
// 	jsonBody, _ := json.Marshal(product)
// 	req := httptest.NewRequest(http.MethodPut, "http://localhost:8080/products/1", bytes.NewBuffer(jsonBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", token)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	c.SetParamNames("id")
// 	c.SetParamValues("1")

// 	if assert.NoError(t, handler.UpdateProduct(c)) {
// 		assert.Equal(t, http.StatusOK, rec.Code)
// 	}
// }

// func TestDeleteProduct(t *testing.T) {
// 	e := echo.New()

// 	req := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/products/1", nil)
// 	req.Header.Set("Authorization", token)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	c.SetParamNames("id")
// 	c.SetParamValues("1")

// 	if assert.NoError(t, handler.DeleteProduct(c)) {
// 		assert.Equal(t, http.StatusOK, rec.Code)
// 	}
// }

// func TestGetProductByID(t *testing.T) {
// 	e := echo.New()

// 	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/products/1", nil)
// 	req.Header.Set("Authorization", token)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	c.SetParamNames("id")
// 	c.SetParamValues("1")

// 	if assert.NoError(t, handler.GetProductByID(c)) {
// 		assert.Equal(t, http.StatusOK, rec.Code)
// 	}
// }

// func TestGetCategoriesWithProducts(t *testing.T) {
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/categories", nil)
// 	req.Header.Set("Authorization", token)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	if assert.NoError(t, handler.GetCategoriesWithProducts(c)) {
// 		assert.Equal(t, http.StatusOK, rec.Code)
// 	}
// }
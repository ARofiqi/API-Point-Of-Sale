package main

import (
	"aro-shop/db"
	"aro-shop/middleware"
	"aro-shop/models"
	"aro-shop/routes"
	"log"

	"github.com/labstack/echo/v4"
)

// "net/http"
// "github.com/labstack/echo/v4/middleware"

func main() {
	e := echo.New()
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"http://127.0.0.1:5500", "http://localhost:3000"}, // Sesuaikan dengan domain frontend
	// 	AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	// }))

	e.Use(middleware.RateLimiterMiddleware)
	
	db.InitDB()
	
	err := db.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Transaction{}, &models.TransactionItem{}, &models.Category{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	routes.SetupRoutes(e)
	e.Start(":8080")
}

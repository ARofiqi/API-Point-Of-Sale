package main

import (
	"aro-shop/db"
	"aro-shop/routes"
	
	"github.com/labstack/echo/v4"
)
// "aro-shop/middleware"
// "net/http"
// "github.com/labstack/echo/v4/middleware"

func main() {
	e := echo.New()
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"http://127.0.0.1:5500", "http://localhost:3000"}, // Sesuaikan dengan domain frontend
	// 	AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	// }))
	
	// e.Use(middleware.RateLimiterMiddleware)

	// db.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})

	db.InitDB()
	routes.SetupRoutes(e)
	e.Start(":8080")
}

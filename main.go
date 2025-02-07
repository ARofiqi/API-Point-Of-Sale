package main

import (
	"aro-shop/db"
	"aro-shop/routes"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://127.0.0.1:5500", "http://localhost:3000"}, // Sesuaikan dengan domain frontend
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	db.InitDB()
	routes.SetupRoutes(e)
	e.Start(":8080")
}

package main

import (
	"aro-shop/config"
	"aro-shop/db"
	"aro-shop/middlewares"
	"aro-shop/queue"
	"aro-shop/routes"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	APPEnv = config.LoadConfig().APPEnv
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://127.0.0.1:5500", "http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	if false {
		e.Use(middlewares.RateLimiterMiddleware)
	}

	db.InitDB()

	queue.InitRabbitMQ()
	defer queue.CloseRabbitMQ()

	go queue.StartWorker()

	routes.SetupRoutes(e)
	e.Start(":8080")
}

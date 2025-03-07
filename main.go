package main

import (
	"aro-shop/config"
	"aro-shop/db"
	"aro-shop/middlewares"
	"aro-shop/routes"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// "aro-shop/rabbitmq"

var (
	APPEnv = config.LoadConfig().APPEnv
)

func main() {
	// err := rabbitmq.InitRabbitMQ("amqp://guest:guest@localhost:5672/")
	// if err != nil {
	// 	log.Fatalf("Tidak dapat terhubung ke RabbitMQ: %s", err)
	// }
	// defer rabbitmq.conn.Close()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://127.0.0.1:5500", "http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	if false {
		e.Use(middlewares.RateLimiterMiddleware)
	}

	db.InitDB()

	routes.SetupRoutes(e)
	e.Start(":8080")
}

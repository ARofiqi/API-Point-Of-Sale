package routes

import (
	"aro-shop/handler"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/products", handler.GetProducts)
	e.POST("/products", handler.CreateProduct)
	e.PUT("/products/:id", handler.UpdateProduct)
	e.DELETE("/products/:id", handler.DeleteProduct)

	e.GET("/transactions", handler.GetTransactions)
	e.POST("/transactions", handler.CreateTransaction)
	e.PUT("/transactions/:id", handler.UpdateTransaction)
	e.DELETE("/transactions/:id", handler.DeleteTransaction)
}

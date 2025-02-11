package routes

import (
	"aro-shop/handler"
	"aro-shop/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.POST("/register", handler.Register)
	e.POST("/login", handler.Login)

	authGroup := e.Group("")
	authGroup.Use(middleware.JWTMiddleware)

	authGroup.GET("/products", handler.GetProducts)
	authGroup.GET("/products/:id", handler.GetProductByID)
	authGroup.POST("/products", handler.CreateProduct)
	authGroup.PUT("/products/:id", handler.UpdateProduct)
	authGroup.DELETE("/products/:id", handler.DeleteProduct)

	authGroup.GET("/categories", handler.GetCategoriesWithProducts)

	authGroup.GET("/transactions", handler.GetTransactions)
	authGroup.POST("/transactions", handler.CreateTransaction)
}

package routes

import (
	"aro-shop/handler"
	"aro-shop/middlewares"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.POST("/register", handler.Register)
	e.POST("/login", handler.Login)

	authGroup := e.Group("")
	authGroup.Use(middlewares.JWTMiddleware)

	authGroup.GET("/products", handler.GetProducts)
	authGroup.GET("/products/:id", handler.GetProductByID)
	authGroup.GET("/categories", handler.GetCategoriesWithProducts)
	authGroup.GET("/transactions", handler.GetTransactions)
	authGroup.GET("/transactions/date", handler.GetTransactionsByDateRange)
	authGroup.GET("/transactions/:id/subtotal", handler.GetTransactionSubtotal)

	adminGroup := e.Group("")
	adminGroup.Use(middlewares.JWTMiddleware, middlewares.RoleMiddleware("admin"))

	adminGroup.POST("/products", handler.CreateProduct)
	adminGroup.PUT("/products/:id", handler.UpdateProduct)
	adminGroup.DELETE("/products	/:id", handler.DeleteProduct)
	adminGroup.POST("/transactions", handler.CreateTransaction)
	adminGroup.PUT("/users/:id/role", handler.SetUserRole)
}

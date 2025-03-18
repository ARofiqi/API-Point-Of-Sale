package routes

import (
	"aro-shop/handler"
	"aro-shop/middlewares"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.POST("/register", handler.Register)
	e.POST("/login", handler.Login)

	authGroup := e.Group("/api")
	authGroup.Use(middlewares.JWTMiddleware)

	authGroup.GET("/products", handler.GetProducts)
	authGroup.GET("/product/:id", handler.GetProductByID)
	authGroup.GET("/category-products", handler.GetCategoriesWithProducts)

	authGroup.GET("/transactions", handler.GetTransactions)
	authGroup.POST("/transaction", handler.CreateTransaction)
	authGroup.GET("/transactions/date", handler.GetTransactionsByDateRange)
	authGroup.GET("/transactions/:id/subtotal", handler.GetTransactionSubtotal)

	authGroup.GET("/notifications", handler.GetNotifications)
	authGroup.PUT("/notification/:id/read", handler.MarkNotificationAsRead)

	authGroup.GET("/categories", handler.GetCategories)
	authGroup.GET("/categories/:id", handler.GetCategory)

	adminGroup := e.Group("/api")
	adminGroup.Use(middlewares.JWTMiddleware, middlewares.RoleMiddleware("admin"))

	adminGroup.PUT("/users/:id/role", handler.SetUserRole)

	adminGroup.POST("/product", handler.CreateProduct)
	adminGroup.PUT("/product/:id", handler.UpdateProduct)
	adminGroup.DELETE("/product/:id", handler.DeleteProduct)

	adminGroup.POST("/categories", handler.CreateCategory)
	adminGroup.PUT("/categories/:id", handler.UpdateCategory)
	adminGroup.DELETE("/categories/:id", handler.DeleteCategory)
}

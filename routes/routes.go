package routes

import (
	"aro-shop/handler"
	"aro-shop/middlewares"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.POST("/api/auth/register", handler.Register)
	e.POST("/api/auth/login", handler.Login)

	superAdminGroup := e.Group("")
	superAdminGroup.Use(middlewares.JWTMiddleware, middlewares.RoleMiddleware("superAdmin"))
	superAdminGroup.POST("/api/auth/register/admin", handler.RegisterAdmin)

	// e.POST("/api/init-superadmin", handler.RegisterAdmin)

	authGroup := e.Group("/api")
	authGroup.Use(middlewares.JWTMiddleware)

	authGroup.GET("/products", handler.GetProducts)
	authGroup.GET("/product/:id", handler.GetProductByID)
	authGroup.GET("/category-products", handler.GetCategoriesWithProducts)

	authGroup.GET("/transactions", handler.GetTransactions)
	authGroup.GET("/transactions/date", handler.GetTransactionsByDateRange)
	authGroup.GET("/transactions/:id/subtotal", handler.GetTransactionSubtotal)
	authGroup.POST("/transaction", handler.CreateTransaction)

	authGroup.GET("/notifications", handler.GetNotifications)
	authGroup.GET("/notification/:id", handler.GetNotificationById)
	authGroup.PUT("/notification/:id/read", handler.MarkNotificationAsRead)

	authGroup.GET("/categories", handler.GetCategories)
	authGroup.GET("/categories/:id", handler.GetCategoriesById)

	authGroup.GET("/paymentMethods", handler.GetPaymentMethods)
	authGroup.GET("/paymentMethods/:id", handler.GetPaymentMethod)
	authGroup.POST("/paymentMethods", handler.CreatePaymentMethod)

	adminGroup := e.Group("/api")
	adminGroup.Use(middlewares.JWTMiddleware, middlewares.RoleMiddleware("admin"))

	adminGroup.POST("/product", handler.CreateProduct)
	adminGroup.PUT("/product/:id", handler.UpdateProduct)
	adminGroup.DELETE("/product/:id", handler.DeleteProduct)

	adminGroup.POST("/categories", handler.CreateCategory)
	adminGroup.PUT("/categories/:id", handler.UpdateCategory)
	adminGroup.DELETE("/categories/:id", handler.DeleteCategory)

	adminGroup.PUT("/paymentMethods/:id", handler.UpdatePaymentMethod)
	adminGroup.DELETE("/paymentMethods/:id", handler.DeletePaymentMethod)
}

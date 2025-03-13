package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/qhh/prjEcom/pkg/api/handlers"
	"github.com/qhh/prjEcom/pkg/api/middlewares"
	"github.com/qhh/prjEcom/pkg/db/store"
	"github.com/qhh/prjEcom/pkg/utils"
)

// SetupRouter sets up all the routes for the API
func SetupRouter(store *store.Store, jwtMaker *utils.JWTMaker) *gin.Engine {
	router := gin.Default()

	// Create handlers
	authHandler := handlers.NewAuthHandler(store, jwtMaker)
	userHandler := handlers.NewUserHandler(store)
	shopHandler := handlers.NewShopHandler(store)
	productHandler := handlers.NewProductHandler(store)
	orderHandler := handlers.NewOrderHandler(store)

	// Auth routes (no authentication required)
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Routes requiring authentication
	api := router.Group("/api")
	api.Use(middlewares.AuthMiddleware(jwtMaker))
	{
		// User routes
		api.GET("/profile", userHandler.GetProfile)

		// Shop routes
		api.POST("/shops", shopHandler.CreateShop)
		api.GET("/shops/user", shopHandler.GetUserShops)
		api.GET("/shops", shopHandler.ListShops)
		api.GET("/shops/search", shopHandler.SearchShops)
		api.GET("/shops/:id", shopHandler.GetShop)
		api.PUT("/shops/:id", shopHandler.UpdateShop)

		// Product routes
		api.GET("/products", productHandler.ListProducts)
		api.GET("/products/search", productHandler.SearchProducts)
		api.GET("/products/category", productHandler.FilterProductsByCategory)
		api.GET("/products/price", productHandler.FilterProductsByPrice)
		api.GET("/products/:id", productHandler.GetProduct)
		api.GET("/shops/:id/products", productHandler.ListProductsByShop)

		// Order routes
		api.POST("/orders", orderHandler.CreateOrder)
		api.GET("/orders/:id", orderHandler.GetOrder)
		api.GET("/orders", orderHandler.GetUserOrders)

		// Seller routes (require seller role)
		seller := api.Group("/seller")
		seller.Use(middlewares.RoleMiddleware("seller", "admin"))
		{
			seller.POST("/products", productHandler.CreateProduct)
			seller.PUT("/products/:id", productHandler.UpdateProduct)
			seller.DELETE("/products/:id", productHandler.DeleteProduct)
		}

		// Admin routes (require admin role)
		admin := api.Group("/admin")
		admin.Use(middlewares.RoleMiddleware("admin"))
		{
			admin.GET("/users", userHandler.ListUsers)
			admin.POST("/users/:id/ban", userHandler.BanUser)
			admin.POST("/users/:id/unban", userHandler.UnbanUser)
			admin.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
		}
	}

	return router
}

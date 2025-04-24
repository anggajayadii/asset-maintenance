package routes

import (
	"asset-maintenance/config"
	"asset-maintenance/controllers"
	"asset-maintenance/middleware"
	"asset-maintenance/repositories"
	"asset-maintenance/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// Auth (Public)
	auth := router.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}

	// Users (Protected - Example only, bisa tambah middleware role check)
	// users := router.Group("/users")
	// users.Use(middleware.JWTAuthMiddleware())
	// {
	// 	users.GET("/", controllers.GetAllUsers)
	// 	users.GET("/:id", controllers.GetUser)
	// 	users.PUT("/:id", controllers.UpdateUser)
	// 	users.DELETE("/:id", controllers.DeleteUser)
	// }

	// Inisialisasi dependencies
	userRepo := repositories.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Group route
	users := router.Group("/users")
	users.Use(middleware.JWTAuthMiddleware()) // Jika perlu auth
	{
		users.GET("/", userController.GetAllUsers)
		users.GET("/:id", userController.GetUser)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
	}

	// // Assets
	// assets := router.Group("/assets")
	// assets.Use(middleware.JWTAuthMiddleware())
	// {
	// 	assets.GET("/", controllers.GetAssets)
	// 	assets.GET("/:id", controllers.GetAssetByID)
	// 	assets.POST("/", controllers.CreateAsset)
	// 	assets.PUT("/:id", controllers.UpdateAsset)
	// 	assets.DELETE("/:id", controllers.DeleteAsset)
	// }

	// Inisialisasi dependencies
	assetRepo := repositories.NewAssetRepository(config.DB)
	assetService := services.NewAssetService(assetRepo)
	assetController := controllers.NewAssetController(assetService)

	// Group route dengan middleware JWT
	assets := router.Group("/assets")
	assets.Use(middleware.JWTAuthMiddleware())
	{
		assets.GET("/", assetController.GetAssets)         // Method dari struct
		assets.GET("/:id", assetController.GetAssetByID)   // Method dari struct
		assets.POST("/", assetController.CreateAsset)      // Method dari struct
		assets.PUT("/:id", assetController.UpdateAsset)    // Method dari struct
		assets.DELETE("/:id", assetController.DeleteAsset) // Method dari struct
	}

	// Maintenance
	// maintenance := router.Group("/maintenance")
	// maintenance.Use(middleware.JWTAuthMiddleware())
	// {
	// 	maintenance.GET("/", controllers.GetMaintenanceRecords)
	// 	maintenance.GET("/:id", controllers.GetMaintenanceByID)
	// 	maintenance.POST("/", controllers.CreateMaintenance)
	// }

	// Inisialisasi dependencies
	maintRepo := repositories.NewMaintenanceRepository(config.DB)
	maintService := services.NewMaintenanceService(maintRepo)
	maintController := controllers.NewMaintenanceController(maintService)

	// Group route
	maintenance := router.Group("/maintenance")
	maintenance.Use(middleware.JWTAuthMiddleware()) // Jika perlu auth
	{
		maintenance.GET("/", maintController.GetRecords)
		maintenance.GET("/:id", maintController.GetRecordByID)
		maintenance.POST("/", maintController.CreateRecord)
	}

	// Asset History
	history := router.Group("/history")
	history.Use(middleware.JWTAuthMiddleware())
	{
		history.GET("/", controllers.GetAssetHistories)
		history.POST("/", controllers.AddAssetHistory)
	}
}

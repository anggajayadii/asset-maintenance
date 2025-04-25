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

	// Users (Admin Only)
	userRepo := repositories.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	users := router.Group("/users")
	users.Use(middleware.JWTAuthMiddleware(), middleware.RoleBasedAuth()) // Tambahkan RBAC
	{
		users.GET("/", userController.GetAllUsers)
		users.GET("/:id", userController.GetUser)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
	}

	// Assets (Role-Specific)
	assetRepo := repositories.NewAssetRepository(config.DB)
	assetService := services.NewAssetService(assetRepo)
	assetController := controllers.NewAssetController(assetService)

	assets := router.Group("/assets")
	assets.Use(middleware.JWTAuthMiddleware(), middleware.RoleBasedAuth()) // Tambahkan RBAC
	{
		assets.GET("/", assetController.GetAssets)
		assets.GET("/:id", assetController.GetAssetByID)
		assets.POST("/", assetController.CreateAsset)
		assets.PUT("/:id", assetController.UpdateAsset)
		assets.DELETE("/:id", assetController.DeleteAsset)
	}

	// Maintenance (Role-Specific)
	maintRepo := repositories.NewMaintenanceRepository(config.DB)
	maintService := services.NewMaintenanceService(maintRepo)
	maintController := controllers.NewMaintenanceController(maintService)

	maintenance := router.Group("/maintenance")
	maintenance.Use(middleware.JWTAuthMiddleware(), middleware.RoleBasedAuth()) // Tambahkan RBAC
	{
		maintenance.GET("/", maintController.GetRecords)
		maintenance.GET("/:id", maintController.GetRecordByID)
		maintenance.POST("/", maintController.CreateRecord)
	}

	// Inisialisasi
	historyRepo := repositories.NewAssetHistoryRepository(config.DB)
	historyService := services.NewAssetHistoryService(historyRepo)
	historyController := controllers.NewAssetHistoryController(historyService)

	// Tambahkan ke route
	history := router.Group("/assets/:asset_id/history")
	history.Use(middleware.JWTAuthMiddleware())
	{
		history.GET("/", historyController.GetAssetHistories)
	}

	// In your router setup:
	historyService := services.NewAssetHistoryService(repositories.NewAssetHistoryRepository(db))
	router.Use(middleware.AssetHistoryMiddleware(historyService))
}

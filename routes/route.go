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
	// Initialize all repositories
	userRepo := repositories.NewUserRepository(config.DB)
	assetRepo := repositories.NewAssetRepository(config.DB)
	maintRepo := repositories.NewMaintenanceRepository(config.DB)
	historyRepo := repositories.NewAssetHistoryRepository(config.DB)

	// Initialize services
	userService := services.NewUserService(userRepo)
	assetService := services.NewAssetService(assetRepo)
	maintService := services.NewMaintenanceService(maintRepo)
	historyService := services.NewAssetHistoryService(historyRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	assetController := controllers.NewAssetController(assetService)
	maintController := controllers.NewMaintenanceController(maintService)
	historyController := controllers.NewAssetHistoryController(historyService)

	// Public routes (no auth required)
	auth := router.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware()) // Apply JWT auth to all protected routes

	// User management (Admin only)
	users := protected.Group("/users")
	users.Use(middleware.RoleBasedAuth())
	{
		users.GET("/", userController.GetAllUsers)
		users.GET("/:id", userController.GetUser)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
	}

	// Asset management
	assets := protected.Group("/assets")
	assets.Use(middleware.RoleBasedAuth())
	{
		assets.GET("/", assetController.GetAssets)
		assets.GET("/:id", assetController.GetAssetByID)
		assets.POST("/", assetController.CreateAsset)
		assets.PUT("/:id", assetController.UpdateAsset)
		assets.DELETE("/:id", assetController.DeleteAsset)
	}

	// Maintenance management
	maintenance := protected.Group("/maintenance")
	maintenance.Use(middleware.RoleBasedAuth())
	{
		maintenance.GET("/", maintController.GetRecords)
		maintenance.GET("/:id", maintController.GetRecordByID)
		maintenance.POST("/", maintController.CreateRecord)
	}

	// Asset history
	history := protected.Group("/asset-history")
	{
		history.GET("/", historyController.GetAssetHistories)
		history.GET("/:asset_id", historyController.GetAssetHistories)
	}
}

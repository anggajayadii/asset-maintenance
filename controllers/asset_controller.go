package controllers

import (
	"asset-maintenance/models"
	"asset-maintenance/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func GetAssets(c *gin.Context) {
// 	var assets []models.Asset
// 	if err := config.DB.Find(&assets).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, assets)
// }

// func GetAssetByID(c *gin.Context) {
// 	id := c.Param("id")
// 	var asset models.Asset
// 	if err := config.DB.First(&asset, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, asset)
// }

// func CreateAsset(c *gin.Context) {
// 	var input models.Asset
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Validasi status
// 	if !input.Status.IsValid() {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Allowed: OK, Dismantle, Defect"})
// 		return
// 	}

// 	userID := c.MustGet("user_id").(uint)
// 	input.AddedBy = userID

// 	if err := config.DB.Create(&input).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, input)
// }

// func UpdateAsset(c *gin.Context) {
// 	id := c.Param("id")
// 	var asset models.Asset

// 	if err := config.DB.First(&asset, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
// 		return
// 	}

// 	var input models.Asset
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Update field yang diizinkan
// 	asset.Name = input.Name
// 	asset.Type = input.Type
// 	asset.DeliveryDate = input.DeliveryDate
// 	asset.Status = input.Status
// 	asset.Location = input.Location
// 	asset.SerialNumber = input.SerialNumber

// 	if err := config.DB.Save(&asset).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, asset)
// }

// func DeleteAsset(c *gin.Context) {
// 	id := c.Param("id")
// 	var asset models.Asset

// 	if err := config.DB.First(&asset, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
// 		return
// 	}

// 	userID := c.MustGet("user_id").(uint)
// 	now := time.Now()

// 	asset.DeletedBy = &userID
// 	asset.DeletedAt = &now

// 	if err := config.DB.Save(&asset).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Soft delete
// 	if err := config.DB.Delete(&asset).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Asset soft-deleted"})
// }

type AssetController struct {
	assetService services.AssetService
}

func NewAssetController(assetService services.AssetService) *AssetController {
	return &AssetController{assetService: assetService}
}

func (ctrl *AssetController) GetAssets(c *gin.Context) {
	assets, err := ctrl.assetService.GetAllAssets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assets)
}

func (ctrl *AssetController) GetAssetByID(c *gin.Context) {
	id := c.Param("id")
	asset, err := ctrl.assetService.GetAssetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}
	c.JSON(http.StatusOK, asset)
}

func (ctrl *AssetController) CreateAsset(c *gin.Context) {
	var input models.Asset
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)
	if err := ctrl.assetService.CreateAsset(&input, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func (ctrl *AssetController) UpdateAsset(c *gin.Context) {
	id := c.Param("id")
	var input models.Asset
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset, err := ctrl.assetService.UpdateAsset(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, asset)
}

func (ctrl *AssetController) DeleteAsset(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(uint)
	if err := ctrl.assetService.DeleteAsset(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Asset soft-deleted"})
}

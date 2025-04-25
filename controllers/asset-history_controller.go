package controllers

import (
	"asset-maintenance/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssetHistoryController struct {
	historyService services.AssetHistoryService
}

func NewAssetHistoryController(historyService services.AssetHistoryService) *AssetHistoryController {
	return &AssetHistoryController{historyService: historyService}
}

func (ctrl *AssetHistoryController) GetAssetHistories(c *gin.Context) {
	assetID, err := strconv.ParseUint(c.Param("asset_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset ID"})
		return
	}

	histories, err := ctrl.historyService.GetAssetHistories(uint(assetID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, histories)
}

package controllers

import (
	"asset-maintenance/config"
	"asset-maintenance/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAssetHistories(c *gin.Context) {
	var histories []models.AssetHistory
	config.DB.Preload("User").Preload("Asset").Find(&histories)
	c.JSON(http.StatusOK, histories)
}

func AddAssetHistory(c *gin.Context) {
	var history models.AssetHistory
	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, history)
}

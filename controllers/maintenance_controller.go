// package controllers

// import (
// 	"asset-maintenance/config"
// 	"asset-maintenance/models"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func GetMaintenanceRecords(c *gin.Context) {
// 	var records []models.Maintenance
// 	config.DB.Preload("Asset").Preload("User").Find(&records)
// 	c.JSON(http.StatusOK, records)
// }

// func GetMaintenanceByID(c *gin.Context) {
// 	id := c.Param("id")
// 	var maint models.Maintenance

// 	if err := config.DB.Preload("Asset").Preload("User").First(&maint, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance record not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, maint)
// }

// func CreateMaintenance(c *gin.Context) {
// 	var maint models.Maintenance
// 	if err := c.ShouldBindJSON(&maint); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if err := config.DB.Create(&maint).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, maint)
// }

package controllers

import (
	"asset-maintenance/models"
	"asset-maintenance/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MaintenanceController struct {
	maintService services.MaintenanceService
}

func NewMaintenanceController(maintService services.MaintenanceService) *MaintenanceController {
	return &MaintenanceController{maintService: maintService}
}

func (ctrl *MaintenanceController) GetRecords(c *gin.Context) {
	records, err := ctrl.maintService.GetAllRecords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (ctrl *MaintenanceController) GetRecordByID(c *gin.Context) {
	id := c.Param("id")
	record, err := ctrl.maintService.GetRecordByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance record not found"})
		return
	}
	c.JSON(http.StatusOK, record)
}

func (ctrl *MaintenanceController) CreateRecord(c *gin.Context) {
	var maint models.Maintenance
	if err := c.ShouldBindJSON(&maint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.maintService.CreateRecord(&maint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, maint)
}

package middleware

import (
	"asset-maintenance/models"
	"asset-maintenance/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AssetHistoryMiddleware(historyService services.AssetHistoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only process for PUT/PATCH requests
		if c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
			// Get the original asset before update
			originalAsset, exists := c.Get("originalAsset")
			if !exists {
				c.Next()
				return
			}

			// Get the updated asset after update
			updatedAsset, exists := c.Get("updatedAsset")
			if !exists {
				c.Next()
				return
			}

			// Get user ID from context (assuming you have auth middleware)
			userID, exists := c.Get("userID")
			if !exists {
				c.Next()
				return
			}

			// Record the changes
			asset, ok1 := originalAsset.(*models.Asset)
			updated, ok2 := updatedAsset.(*models.Asset)
			changedBy, ok3 := userID.(uint)

			if ok1 && ok2 && ok3 {
				err := historyService.RecordAssetUpdate(asset, updated, changedBy)
				if err != nil {
					// Log the error but don't fail the request
					fmt.Printf("Failed to record asset history: %v\n", err)
				}
			}
		}

		c.Next()
	}
}

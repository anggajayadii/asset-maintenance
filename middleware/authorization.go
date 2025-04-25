package middleware

import (
	"asset-maintenance/constants"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RoleBasedAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "role not found"})
			return
		}

		currentPath := c.FullPath()
		currentMethod := c.Request.Method

		// Handle wildcard routes (e.g., /assets/:id)
		basePath := strings.Split(currentPath, "/:")[0]
		if basePath != "" {
			currentPath = basePath
		}

		// Cek permission
		allowed := false
		for path, methods := range constants.RolePermissions[constants.Role(userRole.(string))] {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				for _, m := range methods {
					if m == currentMethod {
						allowed = true
						break
					}
				}
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden,
				gin.H{"error": "your role doesn't have permission"})
			return
		}

		c.Next()
	}
}

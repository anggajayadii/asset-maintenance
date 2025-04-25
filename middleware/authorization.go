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

		// Handle wildcard routes
		basePath := strings.Split(currentPath, "/:")[0]
		if basePath != "" {
			currentPath = basePath
		}

		// Cek permission
		allowed := false
		role := constants.Role(userRole.(string))

		// Dapatkan permissions untuk role ini
		permissions, ok := constants.RolePermissions[role]
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden,
				gin.H{"error": "role permissions not defined"})
			return
		}

		// Cek semua path yang diizinkan
		for path, methods := range permissions {
			if strings.HasPrefix(currentPath, path) {
				for _, method := range methods {
					if method == currentMethod {
						allowed = true
						break
					}
				}
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden,
				gin.H{
					"error":         "your role doesn't have permission",
					"required_role": role,
					"path":          currentPath,
					"method":        currentMethod,
				})
			return
		}

		c.Next()
	}
}

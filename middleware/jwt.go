package middleware

import (
	"asset-maintenance/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from header
		tokenString := extractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required in format: Bearer <token>",
			})
			return
		}

		// Verify and parse the token
		token, err := verifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":  "Invalid token",
				"detail": err.Error(),
			})
			return
		}

		// Check if token is valid and get claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			setAuthContext(c, claims)
			c.Next() // Proceed to the next handler
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			return
		}
	}
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return ""
	}

	return tokenParts[1]
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.JWTSecretKey, nil
	})

	return token, err
}

func setAuthContext(c *gin.Context, claims jwt.MapClaims) {
	// Set user_id to context
	if userID, ok := claims["user_id"].(float64); ok {
		c.Set("user_id", uint(userID))
	} else {
		c.Set("user_id", nil)
	}

	// Set user_role to context
	if role, ok := claims["role"].(string); ok {
		c.Set("user_role", role)
	} else {
		c.Set("user_role", nil)
	}
}

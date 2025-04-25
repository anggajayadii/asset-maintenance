package utils

import (
	"asset-maintenance/config"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID uint
	Role   string // Ubah dari models.Role ke string
	jwt.RegisteredClaims
}

// Ubah parameter role dari string ke models.Role
func GenerateToken(userID uint, role string) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTSecretKey)

	return tokenString, expirationTime, err
}

func ParseToken(tokenStr string) (*Claims, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

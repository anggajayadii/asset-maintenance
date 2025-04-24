package utils

import (
	"asset-maintenance/models"
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
func GenerateToken(userID uint, role models.Role) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	claims := &Claims{
		UserID: userID,
		Role:   string(role), // Konversi ke string
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
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

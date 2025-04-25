package controllers

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"asset-maintenance/config"
	"asset-maintenance/models"
	"asset-maintenance/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	jwtKey        = []byte(os.Getenv("JWT_SECRET_KEY"))
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	passwordRegex = regexp.MustCompile(`^.{8,}$`)
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Username string      `json:"username" binding:"required"`
	Password string      `json:"password" binding:"required"`
	Role     models.Role `json:"role" binding:"required"`
	// Tambahkan field lain yang diperlukan
}

type AuthResponse struct {
	Token     string      `json:"token"`
	User      models.User `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
}

// Register User Baru
func Register(c *gin.Context) {
	var input RegisterInput

	// 1. Bind Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 2. Validasi Username
	input.Username = strings.TrimSpace(input.Username)
	if !usernameRegex.MatchString(input.Username) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username must be 3-20 characters and contain only letters, numbers, and underscore",
		})
		return
	}

	// 3. Validasi Password
	input.Password = strings.TrimSpace(input.Password)
	if !passwordRegex.MatchString(input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 8 characters long",
		})
		return
	}

	// 4. Validasi Role
	if !input.Role.IsValid() {
		validRoles := []string{
			string(models.RoleEngineer),
			string(models.RoleLogistik),
			string(models.RoleManajer),
			string(models.RoleAdmin),
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Invalid role",
			"valid_roles": validRoles,
		})
		return
	}

	// 5. Cek username sudah ada
	var existingUser models.User
	err := config.DB.Where("username = ?", input.Username).First(&existingUser).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// 6. Buat user baru
	newUser := models.User{
		Username: input.Username,
		Role:     input.Role,
		// Tambahkan field lain sesuai kebutuhan
	}

	// 7. Hash Password
	if err := newUser.SetPassword(input.Password); err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 8. Simpan ke Database
	if err := config.DB.Create(&newUser).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 9. Response
	newUser.Password = "" // Jangan kembalikan password
	log.Printf("New user registered: %s (%s)", newUser.Username, newUser.Role)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data": gin.H{
			"user_id":    newUser.UserID,
			"username":   newUser.Username,
			"role":       newUser.Role,
			"created_at": newUser.CreatedAt,
		},
	})
}

// Login dan Generate JWT
func Login(c *gin.Context) {
	var input LoginInput

	// 1. Bind dan validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 2. Bersihkan input
	input.Username = strings.TrimSpace(input.Username)
	input.Password = strings.TrimSpace(input.Password)

	// 3. Validasi input
	if input.Username == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// 4. Cari user di database
	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Delay untuk mencegah timing attack
			time.Sleep(1 * time.Second)
			log.Printf("Login attempt for non-existent user: %s", input.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// 5. Verifikasi password
	if !user.VerifyPassword(input.Password) {
		log.Printf("Failed login attempt for user: %s", input.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 6. Generate token
	tokenString, expiresAt, err := utils.GenerateToken(user.UserID, string(user.Role))
	if err != nil {
		log.Printf("Token generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// // 7. Update last login (optional)
	// user.LastLogin = time.Now()
	// if err := config.DB.Model(&user).Update("last_login", user.LastLogin).Error; err != nil {
	// 	log.Printf("Failed to update last login: %v", err)
	// }

	// 8. Siapkan response
	user.Password = ""
	log.Printf("Successful login: %s (%s)", user.Username, user.Role)

	c.JSON(http.StatusOK, AuthResponse{
		Token:     tokenString,
		User:      user,
		ExpiresAt: expiresAt,
	})

	tokenString, expiresAt, err = utils.GenerateToken(user.UserID, string(user.Role))
	if err != nil {
		log.Printf("Token generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
}

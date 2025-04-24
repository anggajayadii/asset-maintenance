package controllers

import (
	"log"
	"net/http"
	"os"
	"strings"

	"asset-maintenance/config"
	"asset-maintenance/models"
	"asset-maintenance/utils"

	"github.com/gin-gonic/gin"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY")) // Ambil dari environment variable

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Register User Baru
func Register(c *gin.Context) {
	var user models.User

	// 1. Bind Input
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Validasi Manual
	if strings.TrimSpace(user.Username) == "" || strings.TrimSpace(user.Password) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password wajib diisi"})
		return
	}

	// 3. Validasi Role
	switch user.Role {
	case models.RoleEngineer, models.RoleLogistik, models.RoleManajer:
		// Role valid
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role harus engineer, logistik, atau manajer"})
		return
	}

	// 4. Hash Password
	if err := user.SetPassword(user.Password); err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password"})
		return
	}

	// 5. Simpan ke Database
	if err := config.DB.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan user"})
		return
	}

	// 6. Response (tanpa expose password)
	log.Printf("User terdaftar: %s (%s)", user.Username, user.Role)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil",
		"data": gin.H{
			"user_id":  user.UserID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// Login dan Generate JWT
func Login(c *gin.Context) {
	var input LoginInput

	// 1. Bind dan validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid"})
		return
	}

	// 2. Bersihkan input
	input.Username = strings.TrimSpace(input.Username)
	input.Password = strings.TrimSpace(input.Password)

	// 3. Validasi panjang input
	if input.Username == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password tidak boleh kosong"})
		return
	}

	// 4. Cari user di database
	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		log.Printf("Login gagal - User tidak ditemukan: %s", input.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"}) // Pesan umum untuk security
		return
	}

	// 5. Verifikasi password
	if !user.VerifyPassword(input.Password) {
		log.Printf("Login gagal - Password salah untuk user: %s", input.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	// 6. Generate token menggunakan utils yang sudah dibuat
	token, err := utils.GenerateToken(user.UserID, user.Role)
	if err != nil {
		log.Printf("Gagal generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// 7. Siapkan response
	user.Password = "" // Pastikan password tidak dikirim kembali
	log.Printf("Login berhasil: %s (%s)", user.Username, user.Role)

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

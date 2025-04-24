// package controllers

// import (
// 	"asset-maintenance/config"
// 	"asset-maintenance/models"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func GetAllUsers(c *gin.Context) {
// 	var users []models.User
// 	config.DB.Find(&users)
// 	c.JSON(http.StatusOK, users)
// }

// func GetUser(c *gin.Context) {
// 	id := c.Param("id")
// 	var user models.User
// 	if err := config.DB.First(&user, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, user)
// }

// func UpdateUser(c *gin.Context) {
// 	// Ambil ID dari parameter URL
// 	id := c.Param("id")

// 	// Cari user yang akan diupdate
// 	var user models.User
// 	if err := config.DB.First(&user, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	// Bind input data (hanya field yang diizinkan untuk diupdate)
// 	var input struct {
// 		Username string `json:"username"`
// 		FullName string `json:"full_name"`
// 		Role     string `json:"role"`
// 		// Password tidak dimasukkan di sini, sebaiknya buat endpoint terpisah untuk update password
// 	}

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Validasi role jika diupdate
// 	if input.Role != "" {
// 		switch models.Role(input.Role) {
// 		case models.RoleEngineer, models.RoleLogistik, models.RoleManajer:
// 			user.Role = models.Role(input.Role)
// 		default:
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role value"})
// 			return
// 		}
// 	}

// 	// Update field yang diizinkan
// 	if input.Username != "" {
// 		// Cek apakah username sudah digunakan oleh user lain
// 		var existingUser models.User
// 		if err := config.DB.Where("username = ? AND user_id != ?", input.Username, id).First(&existingUser).Error; err == nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
// 			return
// 		}
// 		user.Username = input.Username
// 	}

// 	if input.FullName != "" {
// 		user.FullName = input.FullName
// 	}

// 	// Simpan perubahan
// 	if err := config.DB.Save(&user).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "User updated successfully",
// 		"user":    user,
// 	})
// }

// func DeleteUser(c *gin.Context) {
// 	id := c.Param("id")
// 	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
// }

package controllers

import (
	"asset-maintenance/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	users, err := ctrl.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (ctrl *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := ctrl.userService.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Username string `json:"username"`
		FullName string `json:"full_name"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.userService.UpdateUser(id, services.UpdateUserInput(input))
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "user not found":
			status = http.StatusNotFound
		case "invalid role value", "username already taken":
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.userService.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

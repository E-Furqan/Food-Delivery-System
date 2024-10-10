package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	data "github.com/E-Furqan/Food-Delivery-System/Data"
	config "github.com/E-Furqan/Food-Delivery-System/database_config"
	"github.com/E-Furqan/Food-Delivery-System/utils"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {

	var reg_data data.User

	if err := c.ShouldBindJSON(&reg_data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Where("username = ?", reg_data.Username).First(&reg_data).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exist"})
		return
	}
	// Check if the role exists
	var role data.Role
	if err := config.DB.Where("role_id = ?", reg_data.Role_id).First(&role).Error; err != nil {

		if reg_data.Role_id == "1" {

			role.Role_id = reg_data.Role_id
			role.Role_type = "Customer"

			config.DB.Create(&role)

		} else if reg_data.Role_id == "2" {

			role.Role_id = reg_data.Role_id
			role.Role_type = "Delivery driver"

			config.DB.Create(&role)

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
			return
		}

	}

	if err := config.DB.Create(&reg_data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	var userWithRole data.User
	if err := config.DB.Preload("Role").First(&userWithRole, reg_data.User_id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user with role"})
		return
	}

	c.JSON(http.StatusCreated, userWithRole)
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user data.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.Password != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_at": time.Now().Add(24 * time.Hour).Unix(),
	})

}

func Get_user(c *gin.Context) {

	var user_data []data.User

	if err := config.DB.Preload("Role").Order("User_id asc").Find(&user_data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user_data)

}
func Get_role(c *gin.Context) {
	var user_data []data.Role
	if err := config.DB.Order("Role_id asc").Find(&user_data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user_data)
}

func Update_Role(c *gin.Context) {
	var user data.User

	// Retrieve the username from the context
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Define the input structure for binding
	var input struct {
		Role_id string `json:"role_id" binding:"required"`
	}

	// Bind incoming JSON to input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the Role_id exists in the database
	var role data.Role
	result := config.DB.Where("Role_id = ?", input.Role_id).First(&role)

	if result.Error != nil {
		// Role_id not found in the database
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid Role ID: %v", result.Error)})
		return
	}

	// Fetch the user by username
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update the user's role_id
	user.Role_id = input.Role_id
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully", "role_id": user.Role_id})
}

func Update_user(c *gin.Context) {
	var user data.User

	// Retrieve the username from the context
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	// Perform type assertion to string
	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username type"})
		return
	}

	if err := config.DB.Where("username = ?", usernameStr).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User not found %v %s", err, usernameStr)})
		return
	}

	var update_user data.User
	err1 := c.ShouldBindJSON(&update_user)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}

	if err := config.DB.Model(&user).Updates(update_user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}

	var userWithRole data.User
	if err := config.DB.Preload("Role").First(&userWithRole, user.User_id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user with role"})
		return
	}

	c.JSON(http.StatusCreated, userWithRole)
}

func Delete_user(c *gin.Context) {
	// Retrieve the username from the context
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user data.User
	// Fetch the user by username
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// Delete the user
	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})

}

func Delete_role(c *gin.Context) {

	Admin := os.Getenv("Admin")
	Admin_password := os.Getenv("Admin_pass")
	// Define the input structure for binding
	var input struct {
		Username string `json:"username_admin" binding:"required"`
		Password string `json:"password_admin" binding:"required"`
		Role_id  string `json:"role_id" binding:"required"`
	}

	// Bind incoming JSON to input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var role data.Role
	if input.Username == Admin && input.Password == Admin_password {
		// Fetch the role by role id
		if err := config.DB.Where("Role_id = ?", input.Role_id).First(&role).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		// Delete the role
		if err := config.DB.Delete(&role).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the role"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
	}

}

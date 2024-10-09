package controllers

import (
	"net/http"
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
	// Check if the role exists
	var role data.Role
	if err := config.DB.Where("role_id = ?", reg_data.Role_id).First(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	if err := config.DB.Create(&reg_data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusCreated, reg_data)
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

func Getuser(c *gin.Context) {

	var user_data []data.User

	if err := config.DB.Preload("Role").Order("User_id asc").Find(&user_data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user_data)

}
func Getrole(c *gin.Context) {
	var user_data []data.Role
	if err := config.DB.Order("Role_id asc").Find(&user_data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user_data)

}

func Change_Role(c *gin.Context) {
	var user data.User

	// Retrieve the username from the context
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Define the input structure for binding
	var input struct {
		Role_id string `json:"role_id" binding:"required"` // Corrected from "username" to "role_id"
	}

	// Bind incoming JSON to input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

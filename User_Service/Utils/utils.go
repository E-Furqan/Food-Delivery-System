package utils

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	UserId     uint   `json:"user_id"`
	ActiveRole string `json:"activeRole"`
	jwt.StandardClaims
}

func GenerateResponse(httpStatusCode int, c *gin.Context, title1 string, message1 string, title2 string, input interface{}) {
	response := gin.H{
		title1: message1,
	}

	if title2 != "" && input != nil {
		response[title2] = input
	}

	c.JSON(httpStatusCode, response)
}

func VerifyUserId(c *gin.Context) (any, error) {
	userIdValue, exists := c.Get("userId")
	if !exists {
		return 0, fmt.Errorf("userId does not exist")
	}
	userId, ok := userIdValue.(string)
	if !ok {
		return 0, fmt.Errorf("invalid userId")
	}
	return userId, nil
}
func VerificationRole(c *gin.Context) (any, error) {
	activeRole, exists := c.Get("activeRole")
	if !exists {
		return activeRole, fmt.Errorf("invalid Role Type")
	}

	if activeRole != "Admin" {
		return activeRole, fmt.Errorf("insufficient permissions")
	}

	return activeRole, nil
}

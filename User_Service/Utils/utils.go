package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/gin-gonic/gin"
)

func GenerateResponse(httpStatusCode int, c *gin.Context, title1 string, message1 string, title2 string, input interface{}) {

	errorMessage := strings.TrimPrefix(message1, "ERROR: ")
	response := gin.H{
		title1: errorMessage,
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
	userId, ok := userIdValue.(uint)
	if !ok {
		return 0, fmt.Errorf("invalid userId")
	}
	return userId, nil
}

func VerifyActiveRole(c *gin.Context) (any, error) {
	activeRole, exists := c.Get("activeRole")
	if !exists {
		return activeRole, fmt.Errorf("invalid Role Type")
	}

	if activeRole != "Admin" {
		return activeRole, fmt.Errorf("insufficient permissions")
	}

	return activeRole, nil
}

func CreateUserClaim(user model.User) model.UserClaim {
	var UserClaim model.UserClaim
	UserClaim.UserId = user.UserId
	UserClaim.ActiveRole = user.ActiveRole
	return UserClaim
}

func GetEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func GetAuthToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		GenerateResponse(http.StatusUnauthorized, c, "Message", "authorization token not provided", "error", nil)
		return ""
	}
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		GenerateResponse(http.StatusUnauthorized, c, "Message", "invalid authorization header format", "error", nil)
		return ""
	}
	token := tokenParts[1]
	return token
}
package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Verification(c *gin.Context) (any, error) {
	RestaurantID, exists := c.Get("RestaurantID")
	if !exists {
		return 0, fmt.Errorf("restaurant id does not exist")
	}

	RestaurantID, ok := RestaurantID.(uint)
	if !ok {
		return 0, fmt.Errorf("invalid restaurant id")
	}
	return RestaurantID, nil
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

func GetEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func GetAuthToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization token not provided")
	}
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}
	token := tokenParts[1]
	return token, nil
}

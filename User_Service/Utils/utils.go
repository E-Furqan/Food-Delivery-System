package utils

import (
	"fmt"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey []byte
var refreshTokenKey []byte

type Claims struct {
	Username   string `json:"username"`
	ActiveRole string `json:"activeRole"`
	jwt.StandardClaims
}

func SetEnvValue(envVar environmentVariable.Environment) {
	jwtKey = []byte(envVar.JWT_SECRET)
	refreshTokenKey = []byte(envVar.RefreshTokenKey)
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

func VerificationUsername(c *gin.Context) (any, error) {
	usernameValue, exists := c.Get("username")
	if !exists {
		return 0, fmt.Errorf("username does not exist")
	}
	username, ok := usernameValue.(string)
	if !ok {
		return 0, fmt.Errorf("invalid username")
	}
	return username, nil
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

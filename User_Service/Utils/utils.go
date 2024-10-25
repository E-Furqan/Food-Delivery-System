package utils

import (
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

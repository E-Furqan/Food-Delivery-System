package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GenerateResponse(httpStatusCode int, c *gin.Context, title1 string, message1 string, title2 string, input interface{}) {
	response := gin.H{
		title1: message1,
	}

	if title2 != "" && input != nil {
		response[title2] = input
	}

	c.JSON(httpStatusCode, response)
}

type Claims struct {
	UserId      uint   `json:"user_id"`
	ActiveRole  string `json:"activeRole"`
	ClaimId     uint   `json:"claim_id"`
	ServiceType string `json:"service_type"`
	jwt.StandardClaims
}

package authenticator

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/enviorment_variable"
	"github.com/E-Furqan/Food-Delivery-System/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte(environmentVariable.Get_env("JWT_SECRET"))

// AuthMiddleware is used to protect routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &utils.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token %v %s", err, jwtKey)})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("username", claims.RoleId)
		c.Next()
	}
}

// middle ware
func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := utils.RefreshToken(input.RefreshToken, c)

	if err != nil {
		log.Fatal("Error while refreshing token; ", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expires_at":   time.Now().Add(15 * time.Minute).Unix(), // Adjust based on your access token expiration
	})
}

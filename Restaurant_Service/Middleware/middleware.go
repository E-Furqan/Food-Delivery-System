package Middleware

import (
	"fmt"
	"net/http"
	"strings"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey []byte

func SetEnvValue(envVar environmentVariable.Environment) {
	jwtKey = []byte(envVar.JWT_SECRET)
}

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

			c.JSON(http.StatusUnauthorized, gin.H{"Error": fmt.Sprintf("Invalid token %v %s", err, jwtKey)})
			c.Abort()
			return
		}

		c.Set("RestaurantID", claims.ClaimId)
		c.Next()
	}
}

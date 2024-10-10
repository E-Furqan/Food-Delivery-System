package middleware

import (
	"fmt"
	"net/http"
	"strings"

	environmentvariable "github.com/E-Furqan/Food-Delivery-System/enviorment_variable"
	"github.com/E-Furqan/Food-Delivery-System/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var env = environmentvariable.ReadEnv()
var jwtKey = []byte(env.JWT_SECRET)

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
		c.Next()
	}
}

package Middleware

import (
	"fmt"
	"log"
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
		log.Print(claims.RestaurantID)
		log.Print(claims.ServiceType)
		if err != nil || !token.Valid {

			c.JSON(http.StatusUnauthorized, gin.H{"Error": fmt.Sprintf("Invalid token %v %s", err, jwtKey)})
			c.Abort()
			return
		}

		log.Print("claims.RestaurantId")
		log.Print(claims.RestaurantID)
		c.Set("RestaurantID", claims.RestaurantID)
		c.Next()
	}
}

// func RefreshToken(c *gin.Context) {
// 	var input payload.RefreshToken

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	accessToken, err := utils.RefreshToken(input.RefreshToken, c)

// 	if err != nil {
// 		log.Fatal("Error while refreshing token; ", err)
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"access_token": accessToken,
// 		"expires_at":   time.Now().Add(15 * time.Minute).Unix(),
// 	})
// }
package Middleware

import (
	"fmt"
	"net/http"
	"strings"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (middle *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &model.RestaurantClaim{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(middle.envVar.JWT_SECRET), nil
		})
		if err != nil || !token.Valid {

			c.JSON(http.StatusUnauthorized, gin.H{"Error": fmt.Sprint("Invalid token")})
			c.Abort()
			return
		}

		c.Set("RestaurantID", claims.ClaimId)
		c.Next()
	}
}

func (middle *Middleware) RefreshToken(c *gin.Context) {

	var input model.RefreshToken

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var refreshClaim model.RefreshToken
	refreshClaim.RefreshToken = input.RefreshToken
	refreshClaim.Role = "Restaurant"
	tokens, err := middle.AuthClient.RefreshToken(refreshClaim)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access token":  tokens.AccessToken,
		"refresh token": tokens.RefreshToken,
		"expires at":    tokens.Expiration,
	})
}

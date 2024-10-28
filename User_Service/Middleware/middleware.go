package Middleware

import (
	"fmt"
	"net/http"
	"strings"

	ClientPackage "github.com/E-Furqan/Food-Delivery-System/Client"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	Client *ClientPackage.Client
}

func NewMiddleware(client *ClientPackage.Client) *Middleware {
	return &Middleware{
		Client: client,
	}
}

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

		c.Set("userId", claims.UserId)
		c.Set("activeRole", claims.ActiveRole)
		c.Next()
	}
}

func (ctrl *Middleware) RefreshToken(c *gin.Context) {

	var input payload.RefreshToken

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var refreshClaim payload.RefreshToken
	refreshClaim.RefreshToken = input.RefreshToken
	refreshClaim.ServiceType = "Restaurant"
	tokens, err := ctrl.Client.RefreshToken(refreshClaim)
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

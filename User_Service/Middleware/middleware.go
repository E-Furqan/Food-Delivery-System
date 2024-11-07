package Middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/E-Furqan/Food-Delivery-System/Client/AuthClient"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	AuthClient *AuthClient.AuthClient
	envVar     *model.MiddlewareEnv
}

func NewMiddleware(AuthClient *AuthClient.AuthClient, envVar *model.MiddlewareEnv) *Middleware {
	return &Middleware{
		AuthClient: AuthClient,
		envVar:     envVar,
	}
}

func (middle *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			utils.GenerateResponse(http.StatusNotFound, c, "error", "Authorization token required", "", nil)
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(middle.envVar.JWT_SECRET), nil
		})
		if err != nil || !token.Valid {
			log.Print(middle.envVar.JWT_SECRET)
			utils.GenerateResponse(http.StatusUnauthorized, c, "error", err.Error(), "", nil)
			c.Abort()
			return
		}
		c.Set("userId", claims.ClaimId)
		c.Set("activeRole", claims.Role)
		c.Next()
	}
}

func (ctrl *Middleware) RefreshToken(c *gin.Context) {

	var input model.RefreshToken

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}

	var refreshClaim model.RefreshToken
	refreshClaim.RefreshToken = input.RefreshToken
	refreshClaim.Role = input.Role
	tokens, err := ctrl.AuthClient.RefreshToken(refreshClaim)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Message", "Could not generate token", "error", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access token":  tokens.AccessToken,
		"refresh token": tokens.RefreshToken,
		"expires at":    tokens.Expiration,
	})
}

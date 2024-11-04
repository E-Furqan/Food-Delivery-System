package Middleware

import (
	"net/http"
	"strings"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Middleware struct {
	envVar *environmentVariable.Environment
}

func NewMiddleware(envVar *environmentVariable.Environment) *Middleware {
	return &Middleware{
		envVar: envVar,
	}
}

func (middle *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(middle.envVar.JWT_SECRET), nil
		})
		if err != nil || !token.Valid {

			utils.GenerateResponse(http.StatusUnauthorized, c, "Error", err.Error(), "", nil)
			c.Abort()
			return
		}

		c.Set("activeRole", claims.ActiveRole)
		c.Set("ClaimId", claims.ClaimId)
		c.Next()
	}
}

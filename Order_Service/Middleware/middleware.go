package Middleware

import (
	"net/http"
	"strings"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Middleware struct {
	envVar *model.MiddlewareEnv
}

func NewMiddleware(envVar *model.MiddlewareEnv) *Middleware {
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
		c.Set("ID", claims.ClaimId)
		c.Next()
	}
}

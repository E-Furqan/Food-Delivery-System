package Middleware

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/AuthClient"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	AuthClient AuthClient.AuthClientInterface
	envVar     *model.MiddlewareEnv
}

func NewMiddleware(AuthClient AuthClient.AuthClientInterface, envVar *model.MiddlewareEnv) *Middleware {
	return &Middleware{
		AuthClient: AuthClient,
		envVar:     envVar,
	}
}

type MiddlewareInterface interface {
	AuthMiddleware() gin.HandlerFunc
	RefreshToken(c *gin.Context)
}

package Middleware

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	envVar *model.MiddlewareEnv
}

func NewMiddleware(envVar *model.MiddlewareEnv) *Middleware {
	return &Middleware{
		envVar: envVar,
	}
}

type MiddlewareInterface interface {
	AuthMiddleware() gin.HandlerFunc
}

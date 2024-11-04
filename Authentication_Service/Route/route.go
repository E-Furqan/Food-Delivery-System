package route

import (
	AuthController "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Controller"
	"github.com/gin-gonic/gin"
)

func Auth_routes(server *gin.Engine) {

	auth := server.Group("/auth")
	auth.POST("/generate/token", AuthController.GenerateTokens)
	auth.POST("/refresh/token", AuthController.ReFreshToken)
}

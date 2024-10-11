package route

import (
	controllers "github.com/E-Furqan/Food-Delivery-System/handelers/Controllers"
	"github.com/E-Furqan/Food-Delivery-System/handelers/middleware"
	"github.com/gin-gonic/gin"
)

func User_routes(ctrl *controllers.Controller, server *gin.Engine) {

	server.POST("/user/register", ctrl.Register)
	server.GET("/user/get_users", ctrl.Get_user)
	server.GET("/user/get_role", ctrl.Get_role)
	server.POST("/user/login", ctrl.Login)
	server.POST("/user/refresh_token", ctrl.RefreshToken)

	protected := server.Group("/user")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.PATCH("/update/role", ctrl.Update_Role)
		protected.PATCH("/update/user", ctrl.Update_user)
		protected.DELETE("/delete/user", ctrl.Delete_user)
		protected.DELETE("/delete/role", ctrl.Delete_role)
	}
}

package route

import (
	controllers "github.com/E-Furqan/Food-Delivery-System/Interfaces/Controllers"
	"github.com/E-Furqan/Food-Delivery-System/Interfaces/middleware"
	"github.com/gin-gonic/gin"
)

func User_routes(ctrl *controllers.Controller, server *gin.Engine) {

	server.POST("/app/Register", ctrl.Register)
	server.GET("/app/Get_user", ctrl.Get_user)
	server.GET("/app/Get_role", ctrl.Get_role)
	server.POST("/app/Login", ctrl.Login)

	protected := server.Group("/app")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.PATCH("/update/user/role", ctrl.Update_Role)
		protected.PATCH("/update/user", ctrl.Update_user)
		protected.DELETE("/delete/user", ctrl.Delete_user)
		protected.DELETE("/delete/role", ctrl.Delete_role)
	}
}

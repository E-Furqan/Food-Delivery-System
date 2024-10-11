package route

import (
	authenticator "github.com/E-Furqan/Food-Delivery-System/Authentication"
	UserControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/UserController"
	roleController "github.com/E-Furqan/Food-Delivery-System/Controllers/UserController/RoleControler"

	"github.com/gin-gonic/gin"
)

func User_routes(ctrl *UserControllers.Controller, role *roleController.RoleController, server *gin.Engine) {

	server.POST("/user/register", ctrl.Register)
	server.GET("/user/get_users", ctrl.GetUser)
	server.GET("/user/get_role", role.GetRole)
	server.POST("/user/login", ctrl.Login)
	server.POST("/user/refresh_token", authenticator.RefreshToken)

	protected := server.Group("/user")
	protected.Use(authenticator.AuthMiddleware())
	{
		protected.PATCH("/update/user", ctrl.UpdateUser)
		protected.DELETE("/delete/user", ctrl.DeleteUser)
		protected.DELETE("/delete/role", role.DeleteRole)
	}
}

//grapql

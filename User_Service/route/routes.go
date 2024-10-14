package route

import (
	authenticator "github.com/E-Furqan/Food-Delivery-System/Authentication"
	roleController "github.com/E-Furqan/Food-Delivery-System/Controllers/RoleControler"
	UserControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/UserController"

	"github.com/gin-gonic/gin"
)

func User_routes(ctrl *UserControllers.Controller, role *roleController.RoleController, server *gin.Engine) {

	server.POST("/user/register", ctrl.Register)
	server.POST("/user/login", ctrl.Login)
	server.POST("/user/refresh_token", authenticator.RefreshToken)

	protected := server.Group("/user")
	protected.Use(authenticator.AuthMiddleware())
	{
		protected.GET("/get_role", role.GetRole)
		protected.GET("/get_users", ctrl.GetUser)
		protected.GET("/profile", ctrl.Profile)

		// protected.PATCH("/add/user_roles", role.AddRoleToUser)
		protected.PATCH("/update/profile", ctrl.UpdateUser)

		protected.DELETE("/delete/user", ctrl.DeleteUser)
		protected.DELETE("/delete/role", role.DeleteRole)

		protected.POST("/add/role", role.AddRolesByAdmin)
		protected.POST("/search/user", ctrl.SearchForUser)
	}
}

//grapql

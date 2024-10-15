package route

import (
	authenticator "github.com/E-Furqan/Food-Delivery-System/Authentication"
	roleController "github.com/E-Furqan/Food-Delivery-System/Controllers/RoleControler"
	UserControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/UserController"

	"github.com/gin-gonic/gin"
)

func User_routes(ctrl *UserControllers.Controller, rctrl *roleController.RoleController, server *gin.Engine) {

	server.POST("/user/register", ctrl.Register)
	server.POST("/user/login", ctrl.Login)
	server.POST("/user/refresh_token", authenticator.RefreshToken)

	protected := server.Group("/user")
	protected.Use(authenticator.AuthMiddleware())
	{
		protected.GET("/get_role", rctrl.GetRole)
		protected.GET("/get_users", ctrl.GetUsers)
		protected.GET("/profile", ctrl.Profile)

		// protected.PATCH("/add/user_roles", role.AddRoleToUser)
		protected.PATCH("/update/profile", ctrl.UpdateUser)
		protected.PATCH("/switch/role", ctrl.SwitchRole)

		protected.DELETE("/delete/user", ctrl.DeleteUser)
		protected.DELETE("/delete/role", rctrl.DeleteRole)

		protected.POST("/add/role", rctrl.AddRolesByAdmin)
		protected.POST("/search/user", ctrl.SearchForUser)
	}
}

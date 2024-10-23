package route

import (
	authenticator "github.com/E-Furqan/Food-Delivery-System/Authentication"
	roleController "github.com/E-Furqan/Food-Delivery-System/Controllers/RoleControler"
	UserControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/UserController"

	"github.com/gin-gonic/gin"
)

func User_routes(ctrl *UserControllers.Controller, rCtrl *roleController.RoleController, server *gin.Engine) {

	user := server.Group("/user")
	user.POST("/register", ctrl.Register)
	user.POST("/login", ctrl.Login)
	user.POST("/refresh_token", authenticator.RefreshToken)
	user.POST("/process/order", ctrl.ProcessOrder)

	user.Use(authenticator.AuthMiddleware())
	{
		user.GET("/get_role", rCtrl.GetRole)
		user.GET("/get_users", ctrl.GetUsers)
		user.GET("/profile", ctrl.Profile)

		// protected.PATCH("/add/user_roles", role.AddRoleToUser)
		user.PATCH("/update/profile", ctrl.UpdateUser)
		user.PATCH("/switch/role", ctrl.SwitchRole)

		user.DELETE("/delete/user", ctrl.DeleteUser)
		user.DELETE("/delete/role", rCtrl.DeleteRole)

		user.POST("/add/role", rCtrl.AddRolesByAdmin)
		user.POST("/search/user", ctrl.SearchForUser)
	}
}

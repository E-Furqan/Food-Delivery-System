package route

import (
	roleController "github.com/E-Furqan/Food-Delivery-System/Controllers/RoleControler"
	UserControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/UserController"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"

	"github.com/gin-gonic/gin"
)

func User_routes(ctrl *UserControllers.Controller, rCtrl *roleController.RoleController, middleware *Middleware.Middleware, server *gin.Engine) {

	user := server.Group("/user")
	user.POST("/register", ctrl.Register)
	user.POST("/login", ctrl.Login)
	user.POST("/refresh/token", middleware.RefreshToken)

	user.Use(Middleware.AuthMiddleware())
	{
		user.POST("/process/user/order", ctrl.ProcessOrderUser)
		user.POST("/process/driver/order", ctrl.ProcessOrderDriver)
		user.POST("/assign/driver", ctrl.AssignDriver)

		user.GET("/get/role", rCtrl.GetRole)
		user.GET("/get/users", ctrl.GetUsers)
		user.GET("/profile", ctrl.Profile)
		user.GET("/view/user/orders", ctrl.ViewUserOrders)
		user.GET("/view/driver/orders", ctrl.ViewDriverOrders)
		user.GET("/view/orders/without/driver", ctrl.ViewOrdersWithoutDriver)

		// protected.PATCH("/add/user_roles", role.AddRoleToUser)
		user.PATCH("/update/profile", ctrl.UpdateUser)
		user.PATCH("/switch/role", ctrl.SwitchRole)

		user.DELETE("/delete/user", ctrl.DeleteUser)
		user.DELETE("/delete/role", rCtrl.DeleteRole)

		user.POST("/add/role", rCtrl.AddRolesByAdmin)
		user.POST("/search/user", ctrl.SearchForUser)

	}
}

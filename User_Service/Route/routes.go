package route

import (
	RoleController "github.com/E-Furqan/Food-Delivery-System/Controllers/RoleControler"
	UserControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/UserController"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"

	"github.com/gin-gonic/gin"
)

func User_routes(ctrl UserControllers.UserControllerInterface, rCtrl RoleController.RoleControllerInterface, middleware Middleware.MiddlewareInterface, server *gin.Engine) {

	user := server.Group("/user")
	user.POST("/register", ctrl.Register)
	user.POST("/login", ctrl.Login)
	user.POST("/refresh/token", middleware.RefreshToken)

	user.Use(middleware.AuthMiddleware())
	{
		user.POST("/assign/driver", ctrl.AssignDriver)

		user.GET("/get/roles", rCtrl.GetRoles)
		user.GET("/get/users", ctrl.GetUsers)
		user.GET("/profile", ctrl.Profile)
		user.GET("/view/user/orders", ctrl.ViewUserOrders)
		user.GET("/view/driver/orders", ctrl.ViewDriverOrders)
		user.GET("/view/orders/without/driver", ctrl.ViewOrdersWithoutDriver)

		user.PATCH("/update/order/status", ctrl.UpdateOrderStatus)
		user.PATCH("/update/profile", ctrl.UpdateUser)
		user.PATCH("/switch/role", rCtrl.SwitchRole)

		user.DELETE("/delete/user", ctrl.DeleteUser)
		user.DELETE("/delete/role", rCtrl.DeleteRole)

		user.POST("/add/role", rCtrl.AddRolesByAdmin)
		user.POST("/search/user", ctrl.SearchForUser)

	}
}

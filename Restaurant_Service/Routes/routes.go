package route

import (
	authenticator "github.com/E-Furqan/Food-Delivery-System/Authentication"
	controller "github.com/E-Furqan/Food-Delivery-System/Controllers/RestaurantController"
	"github.com/gin-gonic/gin"
)

func User_routes(ctrl *controller.RestaurantController, server *gin.Engine) {

	restaurantRoute := server.Group("/restaurant")
	restaurantRoute.POST("/register", ctrl.Register)
	restaurantRoute.POST("/login", ctrl.Login)
	restaurantRoute.POST("/refresh/token", authenticator.RefreshToken)
	restaurantRoute.POST("/view/menu", ctrl.ViewMenu)
	restaurantRoute.GET("/get/restaurants", ctrl.GetAllRestaurants)

	restaurantRoute.Use(authenticator.AuthMiddleware())
	{
		restaurantRoute.POST("/add/items", ctrl.AddItemItRestaurantMenu)
		restaurantRoute.DELETE("/delete/items", ctrl.DeleteItemsOfRestaurantMenu)
		restaurantRoute.PATCH("/update/status", ctrl.UpdateRestaurantStatus)
	}
}

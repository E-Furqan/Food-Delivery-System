package route

import (
	Authenticator "github.com/E-Furqan/Food-Delivery-System/Authentication"
	"github.com/E-Furqan/Food-Delivery-System/Controllers/ItemController"
	"github.com/E-Furqan/Food-Delivery-System/Controllers/RestaurantController"
	"github.com/gin-gonic/gin"
)

func User_routes(RestaurantController *RestaurantController.RestaurantController, ItemController *ItemController.ItemController, server *gin.Engine) {

	restaurantRoute := server.Group("/restaurant")
	restaurantRoute.POST("/register", RestaurantController.Register)
	restaurantRoute.POST("/login", RestaurantController.Login)
	restaurantRoute.POST("/refresh/token", Authenticator.RefreshToken)
	restaurantRoute.POST("/view/menu", RestaurantController.ViewMenu)
	restaurantRoute.GET("/get/restaurants", RestaurantController.GetAllRestaurants)
	restaurantRoute.POST("/process/order", RestaurantController.ProcessOrder)
	restaurantRoute.Use(Authenticator.AuthMiddleware())
	{
		restaurantRoute.POST("/add/items", ItemController.AddItemsInRestaurantMenu)
		restaurantRoute.DELETE("/delete/items", ItemController.DeleteItemsOfRestaurantMenu)
		restaurantRoute.DELETE("/cancel/order", RestaurantController.CancelOrder)
		restaurantRoute.PATCH("/update/status", RestaurantController.UpdateRestaurantStatus)
	}
}

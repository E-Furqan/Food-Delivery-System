package route

import (
	"github.com/E-Furqan/Food-Delivery-System/Controllers/ItemController"
	"github.com/E-Furqan/Food-Delivery-System/Controllers/RestaurantController"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"
	"github.com/gin-gonic/gin"
)

func Restaurant_routes(RestaurantController *RestaurantController.RestaurantController, ItemController *ItemController.ItemController, server *gin.Engine) {

	restaurantRoute := server.Group("/restaurant")
	restaurantRoute.POST("/register", RestaurantController.Register)
	restaurantRoute.POST("/login", RestaurantController.Login)
	restaurantRoute.POST("/refresh/token", RestaurantController.RefreshToken)
	restaurantRoute.POST("/view/menu", RestaurantController.ViewMenu)
	restaurantRoute.GET("/get/restaurants", RestaurantController.GetAllRestaurants)
	restaurantRoute.POST("/process/order", RestaurantController.ProcessOrder)
	restaurantRoute.Use(Middleware.AuthMiddleware())
	{
		restaurantRoute.POST("/add/items", ItemController.AddItemsInMenu)
		restaurantRoute.DELETE("/delete/items", ItemController.DeleteItemsFromMenu)
		restaurantRoute.DELETE("/cancel/order", RestaurantController.CancelOrder)
		restaurantRoute.PATCH("/update/status", RestaurantController.UpdateRestaurantStatus)
	}
}

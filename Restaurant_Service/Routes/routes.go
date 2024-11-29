package route

import (
	"github.com/E-Furqan/Food-Delivery-System/Controllers/ItemController"
	"github.com/E-Furqan/Food-Delivery-System/Controllers/RestaurantController"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"
	"github.com/gin-gonic/gin"
)

func Restaurant_routes(RestaurantController RestaurantController.RestaurantControllerInterface,
	ItemController ItemController.ItemControllerInterface,
	middleware Middleware.MiddlewareInterface, server *gin.Engine) {

	restaurantRoute := server.Group("/restaurant")
	restaurantRoute.POST("/register", RestaurantController.Register)
	restaurantRoute.POST("/login", RestaurantController.Login)
	restaurantRoute.POST("/refresh/token", middleware.RefreshToken)
	restaurantRoute.GET("/view/menu", RestaurantController.ViewMenu)
	restaurantRoute.GET("/get/restaurants", RestaurantController.GetAllRestaurants)

	restaurantRoute.Use(middleware.AuthMiddleware())
	{
		restaurantRoute.GET("/get/open/restaurants", RestaurantController.FetchOpenRestaurant)
		restaurantRoute.PATCH("/update/order/status", RestaurantController.UpdateOrderStatus)
		restaurantRoute.POST("/add/items", ItemController.AddItemsInMenu)
		restaurantRoute.DELETE("/delete/items", ItemController.DeleteItemsFromMenu)
		restaurantRoute.PATCH("/update/status", RestaurantController.UpdateRestaurantStatus)
		restaurantRoute.GET("/view/restaurant/orders", RestaurantController.ViewRestaurantOrders)
	}
}

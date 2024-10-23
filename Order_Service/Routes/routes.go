package Routes

import (
	OrderControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/OrderController"
	"github.com/gin-gonic/gin"
)

func User_routes(orderController *OrderControllers.OrderController, server *gin.Engine) {

	orderRoute := server.Group("/order")
	orderRoute.PATCH("/update/status", orderController.UpdateOrderStatus)
	orderRoute.GET("/view/user/orders", orderController.GetOrdersOfUser)
	orderRoute.GET("/view/restaurant/orders", orderController.GetOrdersOfRestaurant)
	orderRoute.POST("/place/order", orderController.PlaceOrder)
	// orderRoute.POST("/checkout", orderController.CheckOut)

}

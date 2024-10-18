package Routes

import (
	OrderControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/OrderController"
	"github.com/gin-gonic/gin"
)

func User_routes(orderController *OrderControllers.OrderController, server *gin.Engine) {

	restaurantRoute := server.Group("/order")
	restaurantRoute.POST("/update/order/status", orderController.UpdateOrderStatus)
	restaurantRoute.POST("/view/order", orderController.GetOrdersOfUser)
	restaurantRoute.POST("/put/order", orderController.PutOrder)

}

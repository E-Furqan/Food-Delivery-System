package Routes

import (
	OrderControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/OrderController"
	"github.com/gin-gonic/gin"
)

func User_routes(orderController *OrderControllers.OrderController, server *gin.Engine) {

	orderRoute := server.Group("/order")
	orderRoute.PATCH("/update/status", orderController.UpdateOrderStatus)
	orderRoute.GET("/view/orders", orderController.GetOrdersOfUser)
	orderRoute.POST("/put/order", orderController.PutOrder)
	orderRoute.POST("/checkout", orderController.CheckOut)

}

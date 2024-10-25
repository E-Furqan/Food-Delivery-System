package Routes

import (
	OrderControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/OrderController"
	"github.com/gin-gonic/gin"
)

func Order_routes(orderController *OrderControllers.OrderController, server *gin.Engine) {

	orderRoute := server.Group("/order")
	orderRoute.PATCH("/update/status", orderController.UpdateOrderStatus)
	orderRoute.GET("/view/user/orders", orderController.GetOrdersOfUser)
	orderRoute.GET("/view/restaurant/orders", orderController.GetOrdersOfRestaurant)
	orderRoute.GET("/generate/invoice", orderController.GenerateInvoice)
	orderRoute.POST("/place/order", orderController.PlaceOrder)
	orderRoute.GET("/view/order", orderController.ViewOrderDetails)
	// orderRoute.POST("/checkout", orderController.CheckOut)

}

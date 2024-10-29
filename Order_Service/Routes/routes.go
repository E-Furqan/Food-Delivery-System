package Routes

import (
	OrderControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/OrderController"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"
	"github.com/gin-gonic/gin"
)

func Order_routes(orderController *OrderControllers.OrderController, middle *Middleware.Middleware, server *gin.Engine) {

	orderRoute := server.Group("/order")
	orderRoute.PATCH("/update/status", orderController.UpdateOrderStatus)
	orderRoute.GET("/view/user/orders", orderController.GetOrdersOfUser)
	orderRoute.GET("/view/restaurant/orders", orderController.GetOrdersOfRestaurant)
	orderRoute.GET("/view/drivers/orders", orderController.GetOrdersOfDeliveryDriver)
	orderRoute.GET("/view/without/drivers/orders", orderController.ViewOrdersWithoutRider)

	orderRoute.Use(middle.AuthMiddleware())
	{
		orderRoute.GET("/generate/invoice", orderController.GenerateInvoice)
		orderRoute.POST("/place/order", orderController.PlaceOrder)
		orderRoute.GET("/view/order", orderController.ViewOrderDetails)
	}
}

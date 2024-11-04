package Routes

import (
	OrderControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/OrderController"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"
	"github.com/gin-gonic/gin"
)

func Order_routes(orderController *OrderControllers.OrderController, middle *Middleware.Middleware, server *gin.Engine) {

	orderRoute := server.Group("/order")

	orderRoute.GET("/view/user/orders", orderController.GetOrdersOfUser)
	orderRoute.GET("/view/restaurant/orders", orderController.GetOrdersOfRestaurant)
	orderRoute.GET("/view/driver/orders", orderController.GetOrdersOfDeliveryDriver)
	orderRoute.GET("/view/without/driver/orders", orderController.ViewOrdersWithoutRider)

	orderRoute.Use(middle.AuthMiddleware())
	{
		orderRoute.PATCH("/update/status", orderController.UpdateOrderStatus)
		orderRoute.PATCH("/assign/diver", orderController.AssignDeliveryDriver)
		orderRoute.GET("/generate/invoice", orderController.GenerateInvoice)
		orderRoute.POST("/place/order", orderController.PlaceOrder)
		orderRoute.GET("/view/order", orderController.ViewOrderDetails)
	}
}

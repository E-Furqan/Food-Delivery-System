package Routes

import (
	OrderControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/OrderController"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"
	_ "github.com/E-Furqan/Food-Delivery-System/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Order_routes(orderController OrderControllers.OrderControllerInterface,
	middle Middleware.MiddlewareInterface, server *gin.Engine) {

	orderRoute := server.Group("/order")
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	orderRoute.Use(middle.AuthMiddleware())
	{

		orderRoute.GET("/view/orders", orderController.GetOrders)
		orderRoute.GET("/view/without/driver/orders", orderController.ViewOrdersWithoutRider)
		orderRoute.GET("/generate/invoice", orderController.GenerateInvoice)
		orderRoute.GET("/view/order", orderController.ViewOrderDetails)

		orderRoute.PATCH("/update/status", orderController.UpdateOrderStatus)
		orderRoute.PATCH("/assign/diver", orderController.AssignDeliveryDriver)

		orderRoute.POST("/place/order", orderController.PlaceOrder)
	}
}

package Routes

import (
	CustomerController "github.com/E-Furqan/Food-Delivery-System/Controllers/CustomerContoller"
	RiderController "github.com/E-Furqan/Food-Delivery-System/Controllers/DeliverRiderController"
	OrderControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/OrderController"
	"github.com/E-Furqan/Food-Delivery-System/Controllers/RestaurantController"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"
	_ "github.com/E-Furqan/Food-Delivery-System/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Order_routes(orderController OrderControllers.OrderControllerInterface,
	restCtrl RestaurantController.RestaurantControllerInterface,
	cusCtrl CustomerController.CustomerControllerInterface,
	riderCtrl RiderController.RiderControllerInterface,
	middle Middleware.MiddlewareInterface, server *gin.Engine) {

	orderRoute := server.Group("/order")
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	orderRoute.GET("/fetch/order/status", orderController.FetchOrderStatus)
	orderRoute.Use(middle.AuthMiddleware())
	{

		orderRoute.GET("/fetch/orders/by/time/frame", orderController.FetchOrdersByTimeFrame)
		orderRoute.GET("/fetch/restaurant/revenue", restCtrl.FetchRevenueOfRestaurants)
		orderRoute.GET("/fetch/top/customers", cusCtrl.FetchTopFiveCustomers)
		orderRoute.GET("/status/frequency", orderController.FetchOrderStatusFrequencies)
		orderRoute.GET("/completed/restaurant/orders", restCtrl.FetchCompletedOrdersCountByRestaurant)
		orderRoute.GET("/top/items", restCtrl.FetchTopPurchasedItems)
		orderRoute.GET("/customer/orders/details", cusCtrl.FetchCustomerOrdersDetails)
		orderRoute.GET("/cancel/orders/details", orderController.FetchCancelOrdersDetails)
		orderRoute.GET("/completed/delivers", riderCtrl.FetchCompletedDeliversRider)
		orderRoute.GET("/Average/order/value", orderController.FetchAverageOrderValue)
		orderRoute.GET("/view/orders", orderController.GetOrders)
		orderRoute.GET("/view/without/driver/orders", orderController.ViewOrdersWithoutRider)
		orderRoute.GET("/generate/invoice", orderController.GenerateInvoice)
		orderRoute.GET("/view/order", orderController.ViewOrderDetails)

		orderRoute.PATCH("/update/status", orderController.UpdateOrderStatus)
		orderRoute.PATCH("/assign/diver", riderCtrl.AssignDeliveryDriver)

		orderRoute.POST("/place/order", orderController.PlaceOrder)
	}
}

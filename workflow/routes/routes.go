package Routes

import (
	"github.com/E-Furqan/Food-Delivery-System/Controllers/orderControllers"
	"github.com/gin-gonic/gin"
)

func Order_routes(orderController orderControllers.OrderControllerInterface, server *gin.Engine) {

	workflow := server.Group("/workflow")
	workflow.GET("/place/order", orderController.PlaceOrder)
}

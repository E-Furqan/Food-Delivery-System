package Routes

import (
	"github.com/E-Furqan/Food-Delivery-System/controllers"
	"github.com/gin-gonic/gin"
)

func Workflow_routes(ctrl controllers.ControllerInterface, server *gin.Engine) {

	workflow := server.Group("/workflow")
	workflow.GET("/place/order", ctrl.PlaceOrder)
	workflow.POST("/datapipeline/sync", ctrl.DataSync)
}

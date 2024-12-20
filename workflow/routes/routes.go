package Routes

import (
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
	"github.com/gin-gonic/gin"
)

func Workflow_routes(workFlow workflows.WorkflowInterface, server *gin.Engine) {

	workflow := server.Group("/workflow")
	workflow.GET("/place/order", workFlow.OrderPlacedWorkflow)
}

package orderControllers

import (
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
	"github.com/gin-gonic/gin"
)

type orderControllers struct {
	WorkFlows workflows.WorkflowInterface
}

func NewController(workFlows workflows.WorkflowInterface) *orderControllers {
	return &orderControllers{
		WorkFlows: workFlows,
	}
}

type OrderControllerInterface interface {
	// ViewUserOrders(c *gin.Context)
	ViewDriverOrders(c *gin.Context)
	PlaceOrder(c *gin.Context)
}

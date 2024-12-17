package orderControllers

import (
	activity "github.com/E-Furqan/Food-Delivery-System/Activity"
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
	"github.com/gin-gonic/gin"
)

type orderControllers struct {
	WorkFlows workflows.WorkflowInterface
	Activity  activity.ActivityInterface
}

func NewController(workFlows workflows.WorkflowInterface, activity activity.ActivityInterface) *orderControllers {
	return &orderControllers{
		WorkFlows: workFlows,
		Activity:  activity,
	}
}

type OrderControllerInterface interface {
	ViewDriverOrders(c *gin.Context)
	PlaceOrder(c *gin.Context)
}

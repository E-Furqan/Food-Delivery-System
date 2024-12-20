package workflows

import (
	activity "github.com/E-Furqan/Food-Delivery-System/Activity"
	"github.com/gin-gonic/gin"
)

type Workflow struct {
	Act activity.ActivityInterface
}

func NewController(act activity.ActivityInterface) *Workflow {
	return &Workflow{
		Act: act,
	}
}

type WorkflowInterface interface {
	OrderPlacedWorkflow(c *gin.Context)
}

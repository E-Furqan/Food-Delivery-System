package workflows

import (
	activity "github.com/E-Furqan/Food-Delivery-System/Activity"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"go.temporal.io/sdk/workflow"
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
	PlaceOrderWorkflow(ctx workflow.Context, order model.CombineOrderItem, token string) error
	DataSyncWorkflow(ctx workflow.Context, pipeline model.Pipeline) error
}

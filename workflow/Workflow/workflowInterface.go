package workflows

import (
	activityPac "github.com/E-Furqan/Food-Delivery-System/Activity"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"go.temporal.io/sdk/workflow"
)

type Workflow struct {
	Act    activityPac.ActivityInterface
	EnvVar model.WorkFlowEnv
}

func NewController(act activityPac.ActivityInterface, envVar model.WorkFlowEnv) *Workflow {
	return &Workflow{
		Act:    act,
		EnvVar: envVar,
	}
}

type WorkflowInterface interface {
	PlaceOrderWorkflow(ctx workflow.Context, order model.CombineOrderItem, token string) error
	DataSyncWorkflow(ctx workflow.Context, pipeline model.Pipeline) error
	MonitorOrderStatus(ctx workflow.Context, createdOrder model.UpdateOrder, token string, email string) error
	MoveDataWorkflow(ctx workflow.Context, sourceToken string, destinationToken string,
		sourceFolderID string, destinationFolderID string,
		sourceConfig model.Config, destinationConfig model.Config,
		batchSize int, counter model.FileCounter, startIndex int) (model.FileCounter, error)
}

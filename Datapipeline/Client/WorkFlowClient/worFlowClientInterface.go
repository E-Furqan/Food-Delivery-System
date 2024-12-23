package workflowClient

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

type WorkflowClient struct {
	envVar model.WorkflowEnv
}

func NewClient(envVar model.WorkflowEnv) *WorkflowClient {
	return &WorkflowClient{
		envVar: envVar,
	}
}

type RestaurantClientInterface interface {
	DatapipelineSync(Pipeline model.Pipeline) error
}

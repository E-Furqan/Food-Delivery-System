package worker

import (
	activityPac "github.com/E-Furqan/Food-Delivery-System/Activity"
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
)

type Worker struct {
	Act      activityPac.ActivityInterface
	WorkFlow workflows.WorkflowInterface
}

func NewController(act activityPac.ActivityInterface, workFlow workflows.WorkflowInterface) *Worker {
	return &Worker{
		Act:      act,
		WorkFlow: workFlow,
	}
}

type WorkerInterface interface {
	WorkerUserStart()
}

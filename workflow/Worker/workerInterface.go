package worker

import (
	activity "github.com/E-Furqan/Food-Delivery-System/Activity"
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
)

type Worker struct {
	Act      activity.ActivityInterface
	WorkFlow workflows.WorkflowInterface
}

func NewController(act activity.ActivityInterface, workFlow workflows.WorkflowInterface) *Worker {
	return &Worker{
		Act:      act,
		WorkFlow: workFlow,
	}
}

type WorkerInterface interface {
	WorkerUserStart()
}

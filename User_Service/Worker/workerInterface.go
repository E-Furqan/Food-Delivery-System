package worker

import (
	activity "github.com/E-Furqan/Food-Delivery-System/Activity"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
)

type Worker struct {
	Repo     database.RepositoryInterface
	Act      activity.ActivityInterface
	WorkFlow workflows.WorkflowInterface
}

func NewController(repo database.RepositoryInterface, act activity.ActivityInterface, workFlow workflows.WorkflowInterface) *Worker {
	return &Worker{
		Repo:     repo,
		Act:      act,
		WorkFlow: workFlow,
	}
}

type WorkerInterface interface {
	WorkerUserStart()
}

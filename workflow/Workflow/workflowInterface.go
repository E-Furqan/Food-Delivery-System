package workflows

import (
	activity "github.com/E-Furqan/Food-Delivery-System/Activity"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"go.temporal.io/sdk/workflow"
)

type Workflow struct {
	Repo database.RepositoryInterface
	Act  activity.ActivityInterface
}

func NewController(repo database.RepositoryInterface, act activity.ActivityInterface) *Workflow {
	return &Workflow{
		Repo: repo,
		Act:  act,
	}
}

type WorkflowInterface interface {
	RegisterWorkflow(ctx workflow.Context, registrationData model.User) error
}

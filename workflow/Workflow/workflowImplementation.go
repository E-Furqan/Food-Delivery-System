package workflows

import (
	"log"
	"time"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (wFlow *Workflow) RegisterWorkflow(ctx workflow.Context, registrationData model.User) error {
	option := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 10,
			MaximumInterval:    time.Second * 30,
			MaximumAttempts:    3,
			BackoffCoefficient: 2.0,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, option)

	err := workflow.ExecuteActivity(ctx, wFlow.Act.RegisterCheckRole, registrationData).Get(ctx, &registrationData)
	if err != nil {
		return err
	}
	log.Print("workflow implementation activity:", registrationData)
	err = workflow.ExecuteActivity(ctx, wFlow.Act.CreateUser, registrationData).Get(ctx, &registrationData)
	if err != nil {
		return err
	}
	log.Print("error user", err)
	return nil
}

func (wFlow *Workflow) ViewDriverOrdersWorkflow(ctx workflow.Context, driverID uint, token string) error {
	option := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 10,
			MaximumInterval:    time.Second * 30,
			MaximumAttempts:    3,
			BackoffCoefficient: 2.0,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, option)

	// var user model.User
	// err := workflow.ExecuteActivity(ctx, wFlow.Act.GetUser, userID).Get(ctx, &user)
	// if err != nil {
	// 	return err
	// }
	var UserOrders []model.UpdateOrder
	log.Print("workflow implementation activity:", driverID)
	err := workflow.ExecuteActivity(ctx, wFlow.Act.ViewOrders, driverID, token).Get(ctx, &UserOrders)
	if err != nil {
		return err
	}

	log.Print("error user", err)
	return nil
}

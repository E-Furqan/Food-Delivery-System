package workflows

import (
	"log"
	"strings"
	"time"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// func (wFlow *Workflow) RegisterWorkflow(ctx workflow.Context, registrationData model.User) error {
// 	option := workflow.ActivityOptions{
// 		StartToCloseTimeout: time.Second * 5,
// 		RetryPolicy: &temporal.RetryPolicy{
// 			InitialInterval:    time.Second * 10,
// 			MaximumInterval:    time.Second * 30,
// 			MaximumAttempts:    3,
// 			BackoffCoefficient: 2.0,
// 		},
// 	}
// 	ctx := workflow.WithActivityOptions(ctx, option)

// 	err := workflow.ExecuteActivity(ctx, wFlow.Act.RegisterCheckRole, registrationData).Get(ctx, &registrationData)
// 	if err != nil {
// 		return err
// 	}
// 	log.Print("workflow implementation activity:", registrationData)
// 	err = workflow.ExecuteActivity(ctx, wFlow.Act.CreateUser, registrationData).Get(ctx, &registrationData)
// 	if err != nil {
// 		return err
// 	}
// 	log.Print("error user", err)
// 	return nil
// }

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

	return nil
}

func (wFlow *Workflow) OrderPlacedWorkflow(ctx workflow.Context, order model.CombineOrderItem, token string) error {
	option := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 10,
			MaximumInterval:    time.Minute * 45,
			MaximumAttempts:    3,
			BackoffCoefficient: 2.0,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, option)

	var items []model.Items
	err := workflow.ExecuteActivity(ctx, wFlow.Act.GetItems, order, token).Get(ctx, &items)
	if err != nil {
		return err
	}
	log.Print("items from get item activity: ", items)

	var totalBill float64
	err = workflow.ExecuteActivity(ctx, wFlow.Act.CalculateBill, order, items).Get(ctx, &totalBill)
	if err != nil {
		return err
	}
	order.TotalBill = totalBill
	log.Print("totalBill from get CalculateBill activity: ", totalBill)

	var createdOrder model.UpdateOrder
	err = workflow.ExecuteActivity(ctx, wFlow.Act.CreateOrder, order, token).Get(ctx, &createdOrder)
	if err != nil {
		return err
	}
	order.TotalBill = totalBill
	log.Print("order returned from order activity: ", createdOrder)

	var message string
	err = workflow.ExecuteActivity(ctx, wFlow.Act.SendEmail, createdOrder.OrderId, createdOrder.OrderStatus, token).Get(ctx, &message)
	if err != nil {
		return err
	}
	log.Print("Email sent successfully: ")

	var delayCounter int
	delayCounter = 0
	for {
		workflow.Sleep(ctx, 2*time.Minute)
		var status string
		err = workflow.ExecuteActivity(ctx, wFlow.Act.CheckOrderStatus, createdOrder.OrderId, token).Get(ctx, &status)
		if err != nil {
			return err
		}
		status = strings.ToLower(status)
		if status == utils.Accepted {
			err = workflow.ExecuteActivity(ctx, wFlow.Act.SendEmail, createdOrder.OrderId, status, token).Get(ctx, &message)
			if err != nil {
				return err
			}
			log.Print("Email sent for accepted order: ", message)
		} else if status == utils.Cancelled {
			err = workflow.ExecuteActivity(ctx, wFlow.Act.SendEmail, createdOrder.OrderId, status, token).Get(ctx, &message)
			if err != nil {
				return err
			}
			log.Print("Email sent for rejected order: ", message)
			break
		} else if status == utils.Completed {
			err = workflow.ExecuteActivity(ctx, wFlow.Act.SendEmail, createdOrder.OrderId, status, token).Get(ctx, &message)
			if err != nil {
				return err
			}
			log.Print("Email sent for completed order: ", message)
			break
		}

		if delayCounter == 1 && status == utils.OrderPlaced {
			err = workflow.ExecuteActivity(ctx, wFlow.Act.SendEmail, createdOrder.OrderId, utils.Delay, token).Get(ctx, &message)
			if err != nil {
				return err
			}
			log.Print("Email sent for delay order: ", message)

		}
		delayCounter += 1
	}

	return nil
}

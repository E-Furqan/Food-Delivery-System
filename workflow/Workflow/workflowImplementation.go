package workflows

import (
	"fmt"
	"log"
	"strings"
	"time"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"go.temporal.io/sdk/workflow"
)

func (wFlow *Workflow) PlaceOrderWorkflow(ctx workflow.Context, order model.CombineOrderItem, token string) error {

	option := utils.ActivityOptions()
	ctx = workflow.WithActivityOptions(ctx, option)
	var message string

	var email model.UserEmail
	err := workflow.ExecuteActivity(ctx, wFlow.Act.FetchUserEmail, token).Get(ctx, &email)
	if err != nil {
		wFlow.SendEmail(ctx, order.UpdateOrder, utils.Cancelled, &message, token, email.Email)
		utils.Sleep(ctx)
		log.Print("error in getting email")
		return err
	}

	var items []model.Items
	err = workflow.ExecuteActivity(ctx, wFlow.Act.GetItems, order, token).Get(ctx, &items)
	if err != nil {
		wFlow.SendEmail(ctx, order.UpdateOrder, utils.Cancelled, &message, token, email.Email)
		utils.Sleep(ctx)
		log.Print("error in get items")
		return err
	}
	log.Print("items from get item activity: ", items)

	var totalBill float64
	err = workflow.ExecuteActivity(ctx, wFlow.Act.CalculateBill, order, items).Get(ctx, &totalBill)
	if err != nil {
		wFlow.SendEmail(ctx, order.UpdateOrder, utils.Cancelled, &message, token, email.Email)
		utils.Sleep(ctx)
		return err
	}
	order.TotalBill = totalBill
	log.Print("totalBill from get CalculateBill activity: ", totalBill)

	var createdOrder model.UpdateOrder
	err = workflow.ExecuteActivity(ctx, wFlow.Act.CreateOrder, order, token).Get(ctx, &createdOrder)
	if err != nil {
		wFlow.SendEmail(ctx, createdOrder, utils.Cancelled, &message, token, email.Email)
		utils.Sleep(ctx)
		return err
	}
	order.TotalBill = totalBill
	log.Print("order returned from order activity: ", createdOrder)

	err = wFlow.SendEmail(ctx, createdOrder, createdOrder.OrderStatus, &message, token, email.Email)
	if err != nil {
		return err
	}
	log.Print("Email sent successfully: ")

	wFlow.DelayOrderChecker(ctx, createdOrder, token, email.Email)

	return nil
}

func (wFlow *Workflow) DelayOrderChecker(ctx workflow.Context, createdOrder model.UpdateOrder, token string, email string) error {
	var delayCounter int
	var message string
	delayCounter = 0

	for {
		workflow.Sleep(ctx, 2*time.Minute)
		var status string
		err := workflow.ExecuteActivity(ctx, wFlow.Act.CheckOrderStatus, createdOrder.OrderId, token).Get(ctx, &status)
		if err != nil {
			return err
		}
		status = strings.ToLower(status)
		if status == utils.Accepted {
			err = wFlow.SendEmail(ctx, createdOrder, status, &message, token, email)
			if err != nil {
				return err
			}
			log.Print("Email sent for accepted order: ", message)
		} else if status == utils.Cancelled {
			err = wFlow.SendEmail(ctx, createdOrder, status, &message, token, email)
			if err != nil {
				return err
			}
			log.Print("Email sent for rejected order: ", message)
			break
		} else if status == utils.Completed {
			err = wFlow.SendEmail(ctx, createdOrder, status, &message, token, email)
			if err != nil {
				return err
			}
			log.Print("Email sent for completed order: ", message)
			break
		}

		if delayCounter == 1 && status == utils.OrderPlaced {
			err = wFlow.SendEmail(ctx, createdOrder, utils.Delay, &message, token, email)
			if err != nil {
				return err
			}
			log.Print("Email sent for delay order: ", message)

		}
		delayCounter += 1

		if delayCounter == 5 && status == utils.OrderPlaced {
			err = wFlow.SendEmail(ctx, createdOrder, utils.Cancelled, &message, token, email)
			if err != nil {
				return err
			}
			log.Print("Email sent for delay order: ", message)
			return fmt.Errorf("order status does not changes within 10 mins")
		}
	}

	return nil
}

func (wFlow *Workflow) SendEmail(ctx workflow.Context, createdOrder model.UpdateOrder, status string, message *string, token string, email string) error {
	err := workflow.ExecuteActivity(ctx, wFlow.Act.SendEmail, createdOrder.OrderId, status, token, email).Get(ctx, &message)
	if err != nil {
		return err
	}
	return nil
}

func (wFlow *Workflow) DataSyncWorkflow(ctx workflow.Context, pipeline model.Pipeline) error {

	option := utils.ActivityOptions()
	ctx = workflow.WithActivityOptions(ctx, option)
	source := utils.CreateSourceObj(pipeline.SourcesID)
	destination := utils.CreateDestinationObj(pipeline.DestinationsID)

	var sourceConfig model.Config
	err := workflow.ExecuteActivity(ctx, wFlow.Act.FetchSourceConfiguration, source).Get(ctx, &sourceConfig)
	if err != nil {
		log.Print("error in fetching source configuration", err.Error())
		return err
	}

	var destinationConfig model.Config
	err = workflow.ExecuteActivity(ctx, wFlow.Act.FetchDestinationConfiguration, destination).Get(ctx, &destinationConfig)
	if err != nil {
		log.Print("error in fetching destination configuration", err.Error())
		return err
	}

	var sourceToken string
	err = workflow.ExecuteActivity(ctx, wFlow.Act.CreateSourceToken, sourceConfig).Get(ctx, &sourceToken)
	if err != nil {
		log.Print("error in creating source client: ", err.Error())
		return err
	}

	var destinationToken string
	err = workflow.ExecuteActivity(ctx, wFlow.Act.CreateDestinationToken, destinationConfig).Get(ctx, &destinationToken)
	if err != nil {
		log.Print("error in creating destination client", err.Error())
		return err
	}

	var counter model.FileCounter
	err = workflow.ExecuteActivity(ctx, wFlow.Act.MoveDataFromSourceToDestination, sourceToken, destinationToken, sourceConfig.FolderURL, destinationConfig.FolderURL, sourceConfig).Get(ctx, &counter)
	if err != nil {
		log.Print("error in fetching moving files", err.Error())
		return err
	}

	err = workflow.ExecuteActivity(ctx, wFlow.Act.AddLogs, counter, pipeline.PipelineID).Get(ctx, &counter)
	if err != nil {
		log.Print("error in fetching source configuration", err.Error())
		return err
	}
	return nil
}

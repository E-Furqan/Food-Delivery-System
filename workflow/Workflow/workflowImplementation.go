package workflows

import (
	"log"
	"net/http"
	"strings"
	"time"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/workflow"
)

func (wFlow *Workflow) OrderPlacedWorkflow(c *gin.Context) {
	var ctx workflow.Context

	var order model.CombineOrderItem
	if err := c.ShouldBindBodyWithJSON(&order); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", "could not bind", "", nil)
		return
	}

	token, err := utils.GetAuthToken(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "message", "could not get token", "error", err)
		return
	}

	option := utils.ActivityOptions()
	ctx = workflow.WithActivityOptions(ctx, option)
	var message string

	var email model.UserEmail
	err = workflow.ExecuteActivity(ctx, wFlow.Act.FetchUserEmail, token).Get(ctx, &email)
	if err != nil {
		wFlow.SendEmail(ctx, order.UpdateOrder, utils.Cancelled, &message, token, email.Email)
		utils.Sleep(ctx)
		log.Print("error in getting email")
		utils.GenerateResponse(http.StatusInternalServerError, c, "message", "could not fetch user email", "error", err)
		return
	}

	var items []model.Items
	err = workflow.ExecuteActivity(ctx, wFlow.Act.GetItems, order, token).Get(ctx, &items)
	if err != nil {
		wFlow.SendEmail(ctx, order.UpdateOrder, utils.Cancelled, &message, token, email.Email)
		utils.Sleep(ctx)
		log.Print("error in get items")
		utils.GenerateResponse(http.StatusInternalServerError, c, "message", "could not get items", "error", err)
		return
	}
	log.Print("items from get item activity: ", items)

	var totalBill float64
	err = workflow.ExecuteActivity(ctx, wFlow.Act.CalculateBill, order, items).Get(ctx, &totalBill)
	if err != nil {
		wFlow.SendEmail(ctx, order.UpdateOrder, utils.Cancelled, &message, token, email.Email)
		utils.Sleep(ctx)
		utils.GenerateResponse(http.StatusInternalServerError, c, "message", "could not calculate bill", "error", err)
		return
	}
	order.TotalBill = totalBill
	log.Print("totalBill from get CalculateBill activity: ", totalBill)

	var createdOrder model.UpdateOrder
	err = workflow.ExecuteActivity(ctx, wFlow.Act.CreateOrder, order, token).Get(ctx, &createdOrder)
	if err != nil {
		wFlow.SendEmail(ctx, createdOrder, utils.Cancelled, &message, token, email.Email)
		utils.Sleep(ctx)
		utils.GenerateResponse(http.StatusInternalServerError, c, "message", "could not place the order", "error", err)
		return
	}
	order.TotalBill = totalBill
	log.Print("order returned from order activity: ", createdOrder)

	err = wFlow.SendEmail(ctx, createdOrder, createdOrder.OrderStatus, &message, token, email.Email)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "message", "could not send email", "error", err)
		return
	}
	log.Print("Email sent successfully: ")

	wFlow.DelayOrderChecker(ctx, createdOrder, token, email.Email)

	utils.GenerateResponse(http.StatusOK, c, "message", "order has been placed successfully", "", nil)
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

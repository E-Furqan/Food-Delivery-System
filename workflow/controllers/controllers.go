package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

func (ctrl *Controller) PlaceOrder(c *gin.Context) {
	var order model.CombineOrderItem
	if err := c.ShouldBindBodyWithJSON(&order); err != nil {
		log.Print("binding issue")
		utils.GenerateResponse(http.StatusBadRequest, c, "error", "could not bind", "", nil)
		return
	}
	log.Print(" issue")
	token, err := utils.GetAuthToken(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "message", "could not get token", "error", err)
		return
	}

	options := client.StartWorkflowOptions{
		ID:                    "place-order-workflow-" + fmt.Sprintf("%v", order.UserID),
		TaskQueue:             model.PlaceOrderTaskQueue,
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
	}

	client_var, err := client.Dial(client.Options{})
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "message", "unable to create Temporal client", "error", err)
		return
	}

	if client_var == nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "message", "Temporal client is nil", "error", err)
		return
	}

	_, err = client_var.ExecuteWorkflow(context.Background(), options, ctrl.WorkFlows.OrderPlacedWorkflow, order, token)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "message", "error in workflow", "error", err)
		return
	}
	c.JSON(http.StatusOK, "fetching drivers order in progress")
}

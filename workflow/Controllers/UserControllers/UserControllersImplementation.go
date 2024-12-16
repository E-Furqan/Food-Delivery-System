package userControllers

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

func (ctrl *Controller) RegisterUser(c *gin.Context) {
	// var registrationData model.User

	// if err := c.ShouldBindJSON(&registrationData); err != nil {
	// 	utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
	// 	return
	// }
	// log.Print("register data:", registrationData)

	// client_var, err := client.Dial(client.Options{})
	// if err != nil {
	// 	utils.GenerateResponse(http.StatusBadRequest, c, "message", "unable to create Temporal client", "error", err)
	// 	return
	// }
	// if client_var == nil {
	// 	utils.GenerateResponse(http.StatusBadRequest, c, "message", "Temporal client is nil", "error", err)
	// 	return
	// }

	// options := client.StartWorkflowOptions{
	// 	ID:                    "registration-workflow-" + registrationData.Username,
	// 	TaskQueue:             model.RegisterTaskQueue,
	// 	WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
	// }

	// _, err = client_var.ExecuteWorkflow(context.Background(), options, ctrl.WorkFlows.RegisterWorkflow, registrationData)
	// if err != nil {
	// 	utils.GenerateResponse(http.StatusBadRequest, c, "message", "error in workflow", "error", err)
	// 	return
	// }
	// log.Print("user created")

	// c.JSON(http.StatusOK, "user registered")
}

func (ctrl *Controller) ViewDriverOrders(c *gin.Context) {

	var userID model.ID
	if err := c.ShouldBindBodyWithJSON(&userID); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", "could not bind", "", nil)
		return
	}

	token, err := utils.GetAuthToken(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "message", "could not get token", "error", err)
		return
	}

	options := client.StartWorkflowOptions{
		ID:                    "view-driver-order-workflow-" + fmt.Sprintf("%v", userID.DeliverDriverID),
		TaskQueue:             model.RegisterTaskQueue,
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

	_, err = client_var.ExecuteWorkflow(context.Background(), options, ctrl.WorkFlows.ViewDriverOrdersWorkflow, userID, token)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "message", "error in workflow", "error", err)
		log.Print("errrir in workflow")
		return
	}

	c.JSON(http.StatusOK, "fetching drivers order in progress")
}

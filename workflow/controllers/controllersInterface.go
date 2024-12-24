package controllers

import (
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	WorkFlows workflows.WorkflowInterface
}

func NewController(workFlows workflows.WorkflowInterface) *Controller {
	return &Controller{
		WorkFlows: workFlows,
	}
}

type ControllerInterface interface {
	PlaceOrder(c *gin.Context)
	DataSync(c *gin.Context)
}

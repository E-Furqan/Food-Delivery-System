package userControllers

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Repo        database.RepositoryInterface
	OrderClient OrderClient.OrdClientInterface
	WorkFlows   workflows.WorkflowInterface
}

func NewController(repo database.RepositoryInterface, OrderClient OrderClient.OrdClientInterface, workFlows workflows.WorkflowInterface) *Controller {
	return &Controller{
		Repo:        repo,
		OrderClient: OrderClient,
		WorkFlows:   workFlows,
	}
}

type UserControllerInterface interface {
	// ViewUserOrders(c *gin.Context)
	ViewUsersOrders(c *gin.Context)
	ViewOrdersWithoutDriver(c *gin.Context)
}

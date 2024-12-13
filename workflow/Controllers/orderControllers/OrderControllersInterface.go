package orderControllers

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/EmailClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
	"github.com/gin-gonic/gin"
)

type orderControllers struct {
	Repo        database.RepositoryInterface
	OrderClient OrderClient.OrdClientInterface
	WorkFlows   workflows.WorkflowInterface
	Email       EmailClient.EmailClientInterface
}

func NewController(repo database.RepositoryInterface, OrderClient OrderClient.OrdClientInterface,
	workFlows workflows.WorkflowInterface, email EmailClient.EmailClientInterface) *orderControllers {
	return &orderControllers{
		Repo:        repo,
		OrderClient: OrderClient,
		WorkFlows:   workFlows,
		Email:       email,
	}
}

type UserControllerInterface interface {
	// ViewUserOrders(c *gin.Context)
	ViewDriverOrders(c *gin.Context)
}

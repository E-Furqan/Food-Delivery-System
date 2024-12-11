package UserControllers

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/AuthClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Repo        database.RepositoryInterface
	OrderClient OrderClient.OrdClientInterface
	AuthClient  AuthClient.AuthClientInterface
	WorkFlows   workflows.WorkflowInterface
}

func NewController(repo database.RepositoryInterface, OrderClient OrderClient.OrdClientInterface, AuthClient AuthClient.AuthClientInterface, workFlows workflows.WorkflowInterface) *Controller {
	return &Controller{
		Repo:        repo,
		OrderClient: OrderClient,
		AuthClient:  AuthClient,
		WorkFlows:   workFlows,
	}
}

type UserControllerInterface interface {
	RegisterWorkflow(c *gin.Context)
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUsers(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	Profile(c *gin.Context)
	SearchForUser(c *gin.Context)
	ViewUserOrders(c *gin.Context)
	UpdateOrderStatus(c *gin.Context)
	ViewDriverOrders(c *gin.Context)
	ViewOrdersWithoutDriver(c *gin.Context)
	AssignDriver(c *gin.Context)
	FetchActiveUser(c *gin.Context)
}

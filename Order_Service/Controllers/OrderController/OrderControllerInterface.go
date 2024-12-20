package OrderControllers

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/RestaurantClient"
	WorkFlow "github.com/E-Furqan/Food-Delivery-System/Client/WorkFlowClient"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	Repo      database.RepositoryInterface
	ResClient RestaurantClient.RestaurantClientInterface
	Workflow  WorkFlow.WorkFlowClientInterface
}

func NewController(repo database.RepositoryInterface, ResClient RestaurantClient.RestaurantClientInterface,
	workflow WorkFlow.WorkFlowClientInterface) *OrderController {
	return &OrderController{
		Repo:      repo,
		ResClient: ResClient,
		Workflow:  workflow,
	}
}

type OrderControllerInterface interface {
	UpdateOrderStatus(c *gin.Context)
	GetOrders(c *gin.Context)
	PlaceOrder(c *gin.Context)
	ViewOrderDetails(c *gin.Context)
	ViewOrdersWithoutRider(c *gin.Context)
	GenerateInvoice(c *gin.Context)
	FetchAverageOrderValue(c *gin.Context)
	FetchCancelOrdersDetails(c *gin.Context)
	FetchOrderStatusFrequencies(c *gin.Context)
	FetchOrdersByTimeFrame(c *gin.Context)
	FetchOrderStatus(c *gin.Context)
	CreateOrder(c *gin.Context)
}

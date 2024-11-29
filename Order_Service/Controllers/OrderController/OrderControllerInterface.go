package OrderControllers

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/RestaurantClient"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	Repo      database.RepositoryInterface
	ResClient RestaurantClient.RestaurantClientInterface
}

func NewController(repo database.RepositoryInterface, ResClient RestaurantClient.RestaurantClientInterface) *OrderController {
	return &OrderController{
		Repo:      repo,
		ResClient: ResClient,
	}
}

type OrderControllerInterface interface {
	UpdateOrderStatus(c *gin.Context)
	AssignDeliveryDriver(c *gin.Context)
	GetOrders(c *gin.Context)
	PlaceOrder(c *gin.Context)
	ViewOrderDetails(c *gin.Context)
	ViewOrdersWithoutRider(c *gin.Context)
	GenerateInvoice(c *gin.Context)
	FetchAverageOrderValue(c *gin.Context)
	FetchCompletedDeliversRider(c *gin.Context)
	FetchCancelOrdersDetails(c *gin.Context)
	FetchCustomerOrdersDetails(c *gin.Context)
	FetchTopPurchasedItems(c *gin.Context)
	FetchCompletedOrdersCountByRestaurant(c *gin.Context)
	FetchOrderStatusFrequencies(c *gin.Context)
}

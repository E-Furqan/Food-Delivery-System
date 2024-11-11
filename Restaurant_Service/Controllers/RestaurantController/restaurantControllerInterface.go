package RestaurantController

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/AuthClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
)

type RestaurantController struct {
	Repo        *database.Repository
	OrderClient *OrderClient.OrderClient
	AuthClient  *AuthClient.AuthClient
}

func NewController(repo *database.Repository, OrderClient *OrderClient.OrderClient, AuthClient *AuthClient.AuthClient) *RestaurantController {
	return &RestaurantController{
		Repo:        repo,
		OrderClient: OrderClient,
		AuthClient:  AuthClient,
	}
}

type RestaurantControllerInterface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetAllRestaurants(c *gin.Context)
	UpdateRestaurantStatus(c *gin.Context)
	ViewMenu(c *gin.Context)
	UpdateOrderStatus(c *gin.Context)
	ViewRestaurantOrders(c *gin.Context)
}

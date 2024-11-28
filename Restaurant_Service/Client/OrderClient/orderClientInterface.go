package OrderClient

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/gin-gonic/gin"
)

type OrderClient struct {
	model.OrderClientEnv
}

func NewClient(env model.OrderClientEnv) *OrderClient {
	return &OrderClient{
		OrderClientEnv: env,
	}
}

type OrdClientInterface interface {
	UpdateOrderStatus(input model.OrderDetails, c *gin.Context) error
	ViewRestaurantOrders(input model.Input, c *gin.Context) (*[]model.OrderDetails, error)
}

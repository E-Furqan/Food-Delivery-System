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
	UpdateOrderStatus(input model.UpdateOrder, c *gin.Context) (*model.UpdateOrder, error)
	AssignDriver(input model.UpdateOrder, c *gin.Context) error
	ViewOrdersWithoutDriver(input model.UpdateOrder, c *gin.Context) (*[]model.UpdateOrder, error)
}

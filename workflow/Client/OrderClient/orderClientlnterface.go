package OrderClient

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
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
	UpdateOrderStatus(input model.UpdateOrder, token string) (*model.UpdateOrder, error)
	CreateOrder(order model.CombineOrderItem, token string) (model.UpdateOrder, error)
	FetchOrderStatus(orderID model.OrderID, token string) (string, error)
}

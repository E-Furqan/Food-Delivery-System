package OrderClient

import model "github.com/E-Furqan/Food-Delivery-System/Models"

type OrderClient struct {
	model.OrderClientEnv
}

func NewClient(env model.OrderClientEnv) *OrderClient {
	return &OrderClient{
		OrderClientEnv: env,
	}
}

type OrdClientInterface interface {
	UpdateOrderStatus(input model.OrderDetails, token string) error
	ViewRestaurantOrders(input model.Input, token string) (*[]model.OrderDetails, error)
	ViewOrdersDetails(input model.OrderDetails, token string) (*model.OrderDetails, error)
}

package RestaurantClient

import model "github.com/E-Furqan/Food-Delivery-System/Models"

type RestaurantClient struct {
	envVar *model.RestaurantClientEnv
}

func NewClient(envVar *model.RestaurantClientEnv) *RestaurantClient {
	return &RestaurantClient{
		envVar: envVar,
	}
}

type RestaurantClientInterface interface {
	GetItems(getItems model.GetItems) ([]model.Items, error)
}

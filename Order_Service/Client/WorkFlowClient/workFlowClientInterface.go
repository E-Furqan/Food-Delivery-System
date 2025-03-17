package WorkFlow

import model "github.com/E-Furqan/Food-Delivery-System/Models"

type WorkFlowClient struct {
	envVar *model.WorkFlowClientEnv
}

func NewClient(envVar *model.WorkFlowClientEnv) *WorkFlowClient {
	return &WorkFlowClient{
		envVar: envVar,
	}
}

type WorkFlowClientInterface interface {
	PlaceORder(order model.CombineOrderItem, token string) error
}

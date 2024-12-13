package EmailClient

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

type EmailClient struct {
	model.OrderClientEnv
}

func NewClient(env model.OrderClientEnv) *EmailClient {
	return &EmailClient{
		OrderClientEnv: env,
	}
}

type EmailClientInterface interface {
	EmailSender(orderID uint, orderStatus string) (string, error)
}

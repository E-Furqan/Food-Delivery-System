package EmailClient

import model "github.com/E-Furqan/Food-Delivery-System/Models"

type EmailClient struct {
	envVar model.EmailEnv
}

func NewClient(envVar model.EmailEnv) *EmailClient {
	return &EmailClient{
		envVar: envVar,
	}
}

type EmailClientInterface interface {
	SendEmail(orderID uint, orderStatus string, userEmail string) (string, error)
}

package userClient

import model "github.com/E-Furqan/Food-Delivery-System/Models"

type UserClient struct {
	envVar model.UserClientEnv
}

func NewClient(envVar model.UserClientEnv) *UserClient {
	return &UserClient{
		envVar: envVar,
	}
}

type UserClientInterface interface {
	FetchEmail(token string) (*model.UserEmail, error)
}

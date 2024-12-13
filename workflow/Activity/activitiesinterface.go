package activity

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

type Activity struct {
}

func NewController() *Activity {
	return &Activity{}
}

type ActivityInterface interface {
	RegisterCheckRole(registrationData model.User) (model.User, error)
	CreateUser(registrationData model.User) (model.User, error)
}

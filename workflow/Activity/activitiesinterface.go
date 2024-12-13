package activity

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
)

type Activity struct {
	Repo database.RepositoryInterface
}

func NewController(repo database.RepositoryInterface) *Activity {
	return &Activity{
		Repo: repo,
	}
}

type ActivityInterface interface {
	RegisterCheckRole(registrationData model.User) (model.User, error)
	CreateUser(registrationData model.User) (model.User, error)
}

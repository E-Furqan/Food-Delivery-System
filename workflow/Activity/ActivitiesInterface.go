package activity

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
)

type Activity struct {
	Repo        database.RepositoryInterface
	OrderClient OrderClient.OrdClientInterface
}

func NewController(repo database.RepositoryInterface, orderClient OrderClient.OrdClientInterface) *Activity {
	return &Activity{
		Repo:        repo,
		OrderClient: orderClient,
	}
}

type ActivityInterface interface {
	RegisterCheckRole(registrationData model.User) (model.User, error)
	CreateUser(registrationData model.User) (model.User, error)
	ViewOrders(UserId uint, token string) (*[]model.UpdateOrder, error)
}

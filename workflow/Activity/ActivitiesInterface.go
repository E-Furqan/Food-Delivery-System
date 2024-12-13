package activity

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/EmailClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
)

type Activity struct {
	Repo        database.RepositoryInterface
	OrderClient OrderClient.OrdClientInterface
	Email       EmailClient.EmailClientInterface
}

func NewController(repo database.RepositoryInterface, orderClient OrderClient.OrdClientInterface, email EmailClient.EmailClientInterface) *Activity {
	return &Activity{
		Repo:        repo,
		OrderClient: orderClient,
		Email:       email,
	}
}

type ActivityInterface interface {
	RegisterCheckRole(registrationData model.User) (model.User, error)
	CreateUser(registrationData model.User) (model.User, error)
	ViewOrders(UserId uint, token string) (*[]model.UpdateOrder, error)
}

package activity

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/EmailClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/RestaurantClient"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
)

type Activity struct {
	Repo        database.RepositoryInterface
	OrderClient OrderClient.OrdClientInterface
	Email       EmailClient.EmailClientInterface
	ResClient   RestaurantClient.RestaurantClientInterface
}

func NewController(repo database.RepositoryInterface, orderClient OrderClient.OrdClientInterface,
	email EmailClient.EmailClientInterface, resClient RestaurantClient.RestaurantClientInterface) *Activity {
	return &Activity{
		Repo:        repo,
		OrderClient: orderClient,
		Email:       email,
		ResClient:   resClient,
	}
}

type ActivityInterface interface {
	RegisterCheckRole(registrationData model.User) (model.User, error)
	CreateUser(registrationData model.User) (model.User, error)
	ViewOrders(UserId uint, token string) (*[]model.UpdateOrder, error)
	UpdateOrderStatus(RestaurantId uint, order model.CombineOrderItem, token string) (string, error)
	SendEmail(orderID uint, orderStatus string, token string) (string, error)
	GetItems(order model.CombineOrderItem, token string) (string, error)
}

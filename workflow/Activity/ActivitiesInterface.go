package activity

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/EmailClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/RestaurantClient"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

type Activity struct {
	OrderClient OrderClient.OrdClientInterface
	Email       EmailClient.EmailClientInterface
	ResClient   RestaurantClient.RestaurantClientInterface
}

func NewController(orderClient OrderClient.OrdClientInterface,
	email EmailClient.EmailClientInterface, resClient RestaurantClient.RestaurantClientInterface) *Activity {
	return &Activity{
		OrderClient: orderClient,
		Email:       email,
		ResClient:   resClient,
	}
}

type ActivityInterface interface {
	GetItems(order model.CombineOrderItem, token string) ([]model.Items, error)
	CalculateBill(CombineOrderItem model.CombineOrderItem, items []model.Items) (float64, error)
	CreateOrder(order model.CombineOrderItem, token string) (model.UpdateOrder, error)
	SendEmail(orderID uint, orderStatus string, token string, userEmail string) (string, error)
	CheckOrderStatus(orderID uint, token string) (string, error)
}

package database

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

type RepositoryInterface interface {
	GetOrders(order *[]model.Order, ID uint, columnName string, SortOrder string, searchColumn string) error
	GetOrder(order *model.Order, OrderId uint) error
	GetOrderWithoutRider(order *[]model.Order) error
	GetOrderItems(orderItems *[]model.OrderItem, orderID uint) error
	GetItemByID(itemID uint, item *model.Item) error
	Update(order *model.Order) error
	PlaceOrder(order *model.Order, CombineOrderItem *model.CombineOrderItem) error
	FetchAllOrder(orders *[]model.Order) error
}

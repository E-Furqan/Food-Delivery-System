package payload

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

type Order struct {
	OrderID     uint   `json:"order_id"`
	UserId      uint   `json:"user_id"`
	OrderStatus string `json:"order_status"`
}

type CombineOrderItem struct {
	Order model.Order
	Items []OrderItemPayload
}

type OrderItemPayload struct {
	ItemId   uint `json:"item_id"`
	Quantity uint `json:"quantity"`
}

type CombineOrderFilter struct {
	Order  model.Order
	Filter Filter
}

type Filter struct {
	ColumnName     string
	OrderDirection string
}

type Items struct {
	ItemId          uint    `gorm:"primaryKey;autoIncrement" json:"item_id"`
	ItemName        string  `gorm:"size:100" json:"item_name"`
	ItemDescription string  `gorm:"size:100" json:"item_description"`
	ItemPrice       float64 `gorm:"type:decimal(10,2)" json:"item_price"`
}

type ProcessOrder struct {
	RestaurantId uint   `json:"restaurant_id"`
	OrderStatus  string `json:"order_status"`
}
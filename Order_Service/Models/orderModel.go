package model

import "time"

type Order struct {
	OrderID          uint      `gorm:"primaryKey;column:order_id;autoIncrement" json:"order_id"`
	UserId           uint      `gorm:"size:100;not null;column:user_id" json:"user_id"`
	RestaurantID     uint      `gorm:"column:restaurant_id" json:"restaurant_id"`
	DeliveryDriverID uint      `gorm:"column:delivery_driver;default:0" json:"delivery_driver"`
	OrderStatus      string    `gorm:"size:100;column:order_status" json:"order_status"`
	TotalBill        float64   `gorm:"size:100;column:total_bill" json:"total_bill"`
	Item             []Item    `gorm:"many2many:order_items;foreignKey:order_id;joinForeignKey:order_id;References:item_id;joinReferences:item_id" json:"items"`
	Time             time.Time `gorm:"autoCreateTime" json:"time"`
}

type OrderItem struct {
	OrderID  uint `gorm:"primaryKey;column:order_id" json:"order_id"`
	ItemId   uint `gorm:"primaryKey;column:item_id" json:"item_id"`
	Quantity uint `gorm:"column:quantity" json:"quantity"`
}

type ID struct {
	OrderID          uint `json:"order_id"`
	RestaurantId     uint `json:"restaurant_id"`
	UserId           uint `json:"user_id"`
	DeliveryDriverID uint `json:"delivery_driver"`
}

type OrderStatusUpdateRequest struct {
	ID
	OrderStatus string `json:"order_status"`
	Role        string `json:"activeRole"`
}

type AssignDeliveryDriver struct {
	OrderID          uint   `json:"order_id"`
	DeliveryDriverID uint   `json:"delivery_driver"`
	Role             string `json:"activeRole"`
}

type CombineOrderItem struct {
	ID
	Items []OrderItemPayload
}

type CombineOrderFilter struct {
	ID
	Filter
}

type ProcessOrder struct {
	ID
	OrderStatus string `json:"order_status"`
}

var UserOrderStatuses = []string{
	"Waiting For Delivery Driver",
	"In for delivery",
	"Delivered",
	"Completed",
}

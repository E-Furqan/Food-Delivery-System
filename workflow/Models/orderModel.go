package model

type CombineOrderItem struct {
	UpdateOrder
	TotalBill float64 `json:"total_bil"`
	Items     []OrderItemPayload
}

type OrderID struct {
	OrderID uint `json:"order_id"`
}

type ID struct {
	OrderId         uint `json:"order_id"`
	UserID          uint `json:"user_id"`
	DeliverDriverID uint `json:"delivery_driver"`
	RestaurantId    uint `json:"restaurant_id"`
}

type UpdateOrder struct {
	ID
	OrderStatus string `json:"order_status"`
}

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
	OrderId      uint `json:"order_id"`
	UserID       uint `json:"user_id"`
	RestaurantId uint `json:"restaurant_id"`
}

type UpdateOrder struct {
	ID
	OrderStatus string `json:"order_status"`
}

type OrderClientEnv struct {
	BASE_URL                string
	UPDATE_ORDER_STATUS_URL string
	ORDER_PORT              string
	Fetch_OrderStatus_URL   string
	Create_Order_URL        string
}

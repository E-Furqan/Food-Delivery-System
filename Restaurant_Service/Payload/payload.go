package payload

type Credentials struct {
	Email    string `json:"restaurant_email"`
	Password string `json:"password"`
}

type SearchOrder struct {
	ColumnName string `json:"column_name"`
	OrderType  string `json:"order_type"`
}

type Input struct {
	ItemId           uint   `json:"item_id"`
	RestaurantId     uint   `json:"restaurant_id"`
	RestaurantStatus string `json:"restaurant_status"`
}

type ProcessOrder struct {
	OrderId      uint   `json:"order_id"`
	RestaurantId uint   `json:"restaurant_id"`
	OrderStatus  string `json:"order_status"`
}

type CombinedInput struct {
	SearchOrder
	Input
}

func GetOrderTransitions() map[string]string {
	return map[string]string{
		"order placed": "Accepted",
		"Accepted":     "In process",
		"In process":   "Waiting For Delivery Driver",
	}
}

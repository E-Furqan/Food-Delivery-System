package payload

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Order struct {
	ColumnName string `json:"column_name"`
	OrderType  string `json:"order_type"`
}

type Input struct {
	ItemId           uint   `json:"item_id"`
	RestaurantId     uint   `json:"restaurant_id"`
	RestaurantStatus string `json:"restaurant_status"`
}

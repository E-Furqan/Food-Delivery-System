package model

type Restaurant struct {
	RestaurantId          uint   `gorm:"primaryKey;autoIncrement" json:"restaurant_id"`
	RestaurantName        string `json:"restaurant_name"`
	RestaurantAddress     string `json:"restaurant_address"`
	RestaurantPhoneNumber string `json:"restaurant_phone_number"`
	RestaurantEmail       string `json:"restaurant_email"`
	Password              string `json:"password"`
	RestaurantStatus      string `json:"restaurant_status"`
	Items                 []Item `json:"items"`
}

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

type CombinedInput struct {
	SearchOrder
	Input
}

type OrderDetails struct {
	OrderID      uint   `json:"order_id"`
	OrderStatus  string `json:"order_status"`
	RestaurantId uint   `json:"restaurant_id"`
}

type OpenRestaurantCount struct {
	OpenRestaurantCount string `json:"open_restaurant_count"`
}

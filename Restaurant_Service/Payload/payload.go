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

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
	ServiceType  string `json:"service_type"`
}
type RestaurantClaim struct {
	ClaimId     uint   `json:"claim_id"`
	ServiceType string `json:"service_type"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiration   int64  `json:"expires_at"`
}

type OrderDetails struct {
	OrderID     uint   `json:"order_id"`
	OrderStatus string `json:"order_status"`
}

func GetOrderTransitions() map[string]string {
	return map[string]string{
		"order placed": "Accepted",
		"Accepted":     "In process",
		"In process":   "Waiting For Delivery Driver",
	}
}

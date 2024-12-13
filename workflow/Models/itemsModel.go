package model

type GetItems struct {
	ColumnName   string `json:"column_name"`
	OrderType    string `json:"order_type"`
	RestaurantId uint   `json:"restaurant_id"`
}

type Items struct {
	ItemId          uint    `json:"item_id"`
	ItemName        string  `json:"item_name"`
	ItemDescription string  `json:"item_description"`
	ItemPrice       float64 `json:"item_price"`
}

type OrderItemPayload struct {
	ItemId   uint `json:"item_id"`
	Quantity uint `json:"quantity"`
}

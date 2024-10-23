package payload

type ID struct {
	OrderID          uint `json:"order_id"`
	RestaurantId     uint `json:"restaurant_id"`
	UserId           uint `json:"user_id"`
	DeliveryDriverID uint `json:"delivery_driver"`
}

type Order struct {
	ID
	OrderStatus string `json:"order_status"`
}

type CombineOrderItem struct {
	ID
	Items []OrderItemPayload
}

type CombineOrderFilter struct {
	ID
	Filter
}

type OrderItemPayload struct {
	ItemId   uint `json:"item_id"`
	Quantity uint `json:"quantity"`
}

type Filter struct {
	ColumnName     string `json:"column_name"`
	OrderDirection string `json:"order_direction"`
}

type Items struct {
	ItemId          uint    `json:"item_id"`
	ItemName        string  `json:"item_name"`
	ItemDescription string  `json:"item_description"`
	ItemPrice       float64 `json:"item_price"`
}

type ProcessOrder struct {
	ID
	OrderStatus string `json:"order_status"`
}

type GetItems struct {
	ColumnName   string `json:"column_name"`
	OrderType    string `json:"order_type"`
	RestaurantId uint   `json:"restaurant_id"`
}

var UserOrderStatuses = []string{
	"Waiting For Delivery Driver",
	"In for delivery",
	"Delivered",
	"Completed",
}

var RestaurantOrderStatuses = []string{
	"Order placed",
	"Accepted",
	"In process",
}

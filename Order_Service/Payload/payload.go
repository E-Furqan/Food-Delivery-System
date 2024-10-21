package payload

import model "github.com/E-Furqan/Food-Delivery-System/Models"

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
	ItemId    uint   `json:"item_id"`
	ItemPrice uint   `json:"item_price"`
	Quantity  uint   `json:"quantity"`
	ItemName  string `json:"item_name"`
}

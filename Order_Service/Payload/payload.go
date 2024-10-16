package payload

type Order struct {
	OrderID     uint   `json:"OrderID"`
	UserId      uint   `json:"UserId"`
	OrderStatus string `json:"OrderStatus"`
}

package model

type CombineOrderItem struct {
	UpdateOrder
	TotalBill float64 `json:"total_bil"`
	Items     []OrderItemPayload
}

type OrderID struct {
	OrderID uint `json:"order_id"`
}

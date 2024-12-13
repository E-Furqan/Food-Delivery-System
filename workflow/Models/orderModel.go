package model

type CombineOrderItem struct {
	UpdateOrder
	Items []OrderItemPayload
}

package model

type ItemID struct {
	ItemId uint `json:"item_id"`
}

type CombinedItemsRestaurantID struct {
	RestaurantId uint     `json:"restaurant_id"`
	Items        []ItemID `json:"items"`
}

type Items struct {
	ItemId          uint    `json:"item_id"`
	ItemName        string  `json:"item_name"`
	ItemDescription string  `json:"item_description"`
	ItemPrice       float64 `json:"item_price"`
	RestaurantId    uint    `json:"restaurant_id"`
}

type OrderItemPayload struct {
	ItemId   uint `json:"item_id"`
	Quantity uint `json:"quantity"`
}

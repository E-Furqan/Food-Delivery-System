package model

type Item struct {
	ItemId          uint    `gorm:"primaryKey;autoIncrement" json:"item_id"`
	ItemName        string  `json:"item_name"`
	ItemDescription string  `json:"item_description"`
	ItemPrice       float64 `json:"item_price"`
	RestaurantId    uint    `json:"restaurant_id"`
}

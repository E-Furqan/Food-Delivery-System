package model

type Item struct {
	ItemId          uint    `gorm:"primaryKey;autoIncrement" json:"item_id"`
	ItemName        string  `gorm:"size:100" json:"item_name"`
	ItemDescription string  `gorm:"size:100" json:"item_description"`
	ItemPrice       float64 `gorm:"type:decimal(10,2)" json:"item_price"`
	RestaurantId    uint    `gorm:"foreignKey" json:"restaurant_id"`
}

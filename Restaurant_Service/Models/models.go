package model

type Restaurant struct {
	RestaurantId          uint   `gorm:"primaryKey;autoIncrement" json:"restaurant_id"`
	RestaurantName        string `gorm:"size:255" json:"restaurant_name"`
	RestaurantAddress     string `gorm:"size:255" json:"restaurant_address"`
	RestaurantPhoneNumber string `gorm:"uniqueIndex;size:100" json:"restaurant_phone_number"`
	RestaurantEmail       string `gorm:"uniqueIndex;size:100" json:"restaurant_email"`
	Password              string `gorm:"size:100" json:"password"`
	RestaurantStatus      string `gorm:"size:50" json:"restaurant_status"`
	Items                 []Item `gorm:"many2many:restaurant_items;" json:"items"`
}

type Item struct {
	ItemId          uint    `gorm:"primaryKey;autoIncrement" json:"item_id"`
	ItemName        string  `gorm:"size:100" json:"item_name"`
	ItemDescription string  `gorm:"size:100" json:"item_description"`
	ItemPrice       float64 `gorm:"type:decimal(10,2)" json:"item_price"`
}

type RestaurantItem struct {
	RestaurantId uint `gorm:"primaryKey" json:"restaurant_id"`
	ItemId       uint `gorm:"primaryKey" json:"item_id"`
}

package model

type Restaurant struct {
	RestaurantId          uint   `gorm:"primaryKey;autoIncrement" json:"RestaurantId"`
	RestaurantName        string `gorm:"size:255" json:"RestaurantName"`
	RestaurantAddress     string `gorm:"size:255" json:"RestaurantAddress"`
	RestaurantPhoneNumber string `gorm:"size:100" json:"RestaurantPhoneNumber"`
	RestaurantEmail       string `gorm:"size:100" json:"RestaurantEmail"`
	RestaurantStatus      string `gorm:"size:50" json:"RestaurantStatus"`
	Items                 []Item `gorm:"many2many:restaurant_items;" json:"Items"`
}

type Item struct {
	ItemId   uint   `gorm:"primaryKey;autoIncrement" json:"ItemId"`
	ItemName string `gorm:"size:100" json:"ItemName"`
	RoleType string `gorm:"size:100" json:"RoleType"`
}

type RestaurantItem struct {
	RestaurantId uint `gorm:"primaryKey" json:"RestaurantId"`
	ItemId       uint `gorm:"primaryKey" json:"ItemId"`
}

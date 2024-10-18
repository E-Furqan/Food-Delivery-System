package model

type Order struct {
	OrderID      uint   `gorm:"primaryKey;column:order_id;autoIncrement" json:"order_id"`
	UserId       uint   `gorm:"size:100;not null;column:user_id" json:"user_id"`
	RestaurantID uint   `gorm:"column:restaurant_id" json:"restaurant_id"`
	OrderStatus  string `gorm:"size:100;column:order_status" json:"order_status"`
	TotalBill    uint   `gorm:"size:100;column:total_bill" json:"total_bill"`
	Item         []Item `gorm:"many2many:order_items;foreignKey:order_id;joinForeignKey:order_id;References:item_id;joinReferences:item_id" json:"items"`
}

type Item struct {
	ItemId   uint   `gorm:"primaryKey;size:100;not null;column:item_id" json:"item_id"`
	ItemName string `gorm:"size:100;not null;column:item_name" json:"item_name"`
}

type OrderItem struct {
	OrderID      uint `gorm:"primaryKey;column:order_id" json:"order_id"`
	ItemId       uint `gorm:"primaryKey;column:item_id" json:"item_id"`
	RestaurantID uint `gorm:"primaryKey;column:restaurant_id" json:"restaurant_id"`
	ItemPrice    uint `gorm:"column:item_price" json:"item_price"`
	Quantity     uint `gorm:"column:quantity" json:"quantity"`
}

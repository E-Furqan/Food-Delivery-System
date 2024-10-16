package model

type Order struct {
	OrderID      uint   `gorm:"primaryKey;column:OrderId;autoIncrement" json:"OrderId"`
	UserId       uint   `gorm:"size:100;not null;column:UserId" json:"UserId"`
	RestaurantID uint   `gorm:"column:RestaurantId" json:"RestaurantId"`
	OrderStatus  string `gorm:"size:100;column:OrderStatus" json:"OrderStatus"`
	TotalBill    uint   `gorm:"size:100;column:TotalBill" json:"TotalBill"`
}

type OrderItem struct {
	OrderID   uint `gorm:"column:OrderId" json:"OrderId"`
	ItemId    uint `gorm:"size:100;not null;column:ItemId" json:"ItemId"`
	Quantity  uint `gorm:"column:Quantity" json:"Quantity"`
	ItemPrice uint `gorm:"column:ItemPrice" json:"ItemPrice"`
}

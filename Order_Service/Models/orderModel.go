package model

import "time"

type Order struct {
	OrderID        uint      `gorm:"primaryKey;column:order_id;autoIncrement" json:"order_id"`
	UserId         uint      `gorm:"size:100;not null;column:user_id" json:"user_id"`
	RestaurantID   uint      `gorm:"column:restaurant_id" json:"restaurant_id"`
	DeliveryDriver uint      `gorm:"column:delivery_driver;default:0" json:"delivery_driver"`
	OrderStatus    string    `gorm:"size:100;column:order_status" json:"order_status"`
	TotalBill      float64   `gorm:"size:100;column:total_bill" json:"total_bill"`
	Item           []Item    `gorm:"many2many:order_items;foreignKey:order_id;joinForeignKey:order_id;References:item_id;joinReferences:item_id" json:"items"`
	Time           time.Time `gorm:"autoCreateTime" json:"time"`
}

type OrderItem struct {
	OrderID  uint `gorm:"primaryKey;column:order_id" json:"order_id"`
	ItemId   uint `gorm:"primaryKey;column:item_id" json:"item_id"`
	Quantity uint `gorm:"column:quantity" json:"quantity"`
}

type ID struct {
	OrderID        uint `json:"order_id"`
	RestaurantId   uint `json:"restaurant_id"`
	UserId         uint `json:"user_id"`
	DeliveryDriver uint `json:"delivery_driver"`
}

type OrderStatusUpdateRequest struct {
	ID
	OrderStatus string `json:"order_status"`
	Role        string `json:"activeRole"`
}

type AssignDeliveryDriver struct {
	OrderID        uint   `json:"order_id"`
	DeliveryDriver uint   `json:"delivery_driver"`
	Role           string `json:"activeRole"`
}

type CombineOrderItem struct {
	ID
	TotalBill float64 `json:"total_bill"`
	Items     []OrderItemPayload
}

type ProcessOrder struct {
	ID
	OrderStatus string `json:"order_status"`
}

// AverageOrderValue godoc
// @Description Model for average order value input
type AverageOrderValue struct {
	TimeRange
	FilterType string `json:"filter_type"`
}

// UserAverageOrderValue godoc
// @Description Model for average user order value output
type UserAverageOrderValue struct {
	UserId            uint    `json:"user_id"`
	AverageOrderValue float64 `json:"average_order_value"`
}

// RestaurantAverageOrderValue godoc
// @Description Model for average restaurant order value output
type RestaurantAverageOrderValue struct {
	RestaurantId      uint    `json:"restaurant_id"`
	AverageOrderValue float64 `json:"average_order_value"`
}

// TimeAverageOrderValue godoc
// @Description Model for average order value within a time range output
type TimeAverageOrderValue struct {
	Time              time.Time `json:"time"`
	AverageOrderValue float64   `json:"average_order_value"`
}

// CompletedDelivers godoc
// @Description Model for completed delivers output
type CompletedDelivers struct {
	DeliveryDriver    uint `json:"delivery_driver"`
	CompletedDelivers int  `json:"completed_delivers"`
}

// OrderDetails godoc
// @Description Model for canceled order details output
type OrderDetails struct {
	OrderID        uint      `json:"order_id"`
	UserId         uint      `json:"user_id"`
	RestaurantID   uint      `json:"restaurant_id"`
	ItemDetails    string    `json:"item_details"`
	TotalBill      float64   `json:"total_bill"`
	DeliveryDriver uint      `json:"delivery_driver"`
	OrderStatus    string    `json:"order_status"`
	OrderTime      time.Time `json:"order_time"`
}

type PageNumber struct {
	PageNumber int `json:"page_number" binding:"required"`
	Limit      int `json:"limit" binding:"required"`
}

type MostPurchasedItem struct {
	RestaurantID  uint `json:"restaurant_id"`
	ItemID        uint `json:"item_id"`
	PurchaseCount int  `json:"purchase_count"`
}

type RestaurantCompletedOrdersCount struct {
	RestaurantID    uint `json:"restaurant_id"`
	CompletedOrders int  `json:"completed_orders"`
}

type OrderStatusFrequency struct {
	OrderStatus     string `json:"order_status"`
	StatusFrequency int    `json:"status_frequency"`
}

type UserOrderFrequency struct {
	UserId         uint `json:"user_id"`
	OrderFrequency int  `json:"order_frequency"`
}

type RestaurantRevenue struct {
	RestaurantID uint    `json:"restaurant_id"`
	Revenue      float64 `json:"revenue"`
}
type TimeRange struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type TimeFrame struct {
	TimeFrame string `json:"time_frame" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type OrdersByDay struct {
	HoursOfDate string `json:"hours_of_date"`
	TotalOrders int    `json:"total_orders"`
}

type OrdersByWeek struct {
	DateOfWeek  string `json:"date_of_week"`
	TotalOrders int    `json:"total_orders"`
}

type OrderByMonth struct {
	WeekOfMonth string `json:"week_of_month"`
	TotalOrders int    `json:"total_orders"`
}

type OrderByYear struct {
	MonthOfYear string `json:"month_of_year"`
	TotalOrders int    `json:"total_orders"`
}

var UserOrderStatuses = []string{
	"Waiting For Delivery Driver",
	"In for delivery",
	"Delivered",
	"Completed",
}

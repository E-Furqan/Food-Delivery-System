package database

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

type RepositoryInterface interface {
	GetOrders(order *[]model.Order, ID uint, columnName string, SortOrder string, searchColumn string) error
	GetOrder(order *model.Order, OrderId uint) error
	GetOrderWithoutRider(order *[]model.Order) error
	GetOrderItems(orderItems *[]model.OrderItem, orderID uint) error
	GetItemByID(itemID uint, item *model.Item) error
	Update(order *model.Order) error
	PlaceOrder(order *model.Order, CombineOrderItem *model.CombineOrderItem) error
	FetchAllOrder(orders *[]model.Order) error
	FetchAverageOrderValueOfUser(input uint) (model.UserAverageOrderValue, error)
	FetchAverageOrderValueOfRestaurant(input uint) (model.RestaurantAverageOrderValue, error)
	FetchAverageOrderValueBetweenTime(startTime string, endTime string) (model.TimeAverageOrderValue, error)
	FetchCompletedDeliversOfRider() ([]model.CompletedDelivers, error)
	FetchCancelledOrdersWithItemDetails(Limit int, offset int) ([]model.OrderDetails, error)
	FetchUserOrdersWithItemDetails(userID, limit, offset int) ([]model.OrderDetails, error)
	FetchTopPurchasedItems() ([]model.MostPurchasedItem, error)
	FetchCompletedOrdersCountByRestaurant(timeRange model.TimeRange) ([]model.RestaurantCompletedOrdersCount, error)
	FetchOrderStatusFrequencies() ([]model.OrderStatusFrequency, error)
	FetchTopFiveCustomers() ([]model.UserOrderFrequency, error)
	FetchRestaurantsRevenue() ([]model.RestaurantRevenue, error)
	FetchOrdersByDay(timeRange model.TimeFrame) ([]model.OrdersByDay, error)
	FetchOrdersByWeek(timeRange model.TimeFrame) ([]model.OrdersByWeek, error)
	FetchOrdersByMonth(timeRange model.TimeFrame) ([]model.OrderByMonth, error)
	FetchOrdersByYear(timeRange model.TimeFrame) ([]model.OrderByYear, error)
	FetchOrderStatus(orderId uint) (model.ProcessOrder, error)
}

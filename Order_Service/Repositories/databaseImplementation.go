package database

import (
	"fmt"
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"

	"gorm.io/gorm"
)

func (repo *Repository) GetOrders(order *[]model.Order, ID uint, columnName string, SortOrder string, searchColumn string) error {

	if SortOrder != "asc" && SortOrder != "desc" {
		SortOrder = "asc"
	}

	validColumns := map[string]bool{
		"order_id":      true,
		"user_id":       true,
		"restaurant_id": true,
		"order_status":  true,
	}

	if !validColumns[columnName] {
		columnName = "order_id"
	}

	if !validColumns[searchColumn] {
		return fmt.Errorf("invalid search column: %s", searchColumn)
	}

	tx := repo.DB.Begin()
	err := repo.DB.Where((fmt.Sprintf("%s = ?", searchColumn)), ID).Preload("Item").Order(fmt.Sprintf("%s %s", columnName, SortOrder)).Find(order).Error
	if err != nil {
		tx.Rollback()
		return nil
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (repo *Repository) GetOrder(order *model.Order, OrderId uint) error {
	tx := repo.DB.Begin()
	err := tx.Where("order_id = ?", OrderId).Preload("Item").First(order).Error
	log.Printf("errror: %s", err)
	log.Print("order: ", order)
	log.Printf("OrderId: %v", OrderId)
	if err != nil {
		tx.Rollback()
		return nil
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (repo *Repository) GetOrderWithoutRider(order *[]model.Order) error {
	tx := repo.DB.Begin()
	err := repo.DB.Where("delivery_driver = ?", 0).Preload("Item").Find(order).Error
	if err != nil {
		tx.Rollback()
		return nil
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (repo *Repository) GetOrderItems(orderItems *[]model.OrderItem, orderID uint) error {
	tx := repo.DB.Begin()
	err := repo.DB.Where("order_id = ?", orderID).Find(orderItems).Error
	if err != nil {
		tx.Rollback()
		return nil
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (repo *Repository) GetItemByID(itemID uint, item *model.Item) error {
	tx := repo.DB.Begin()
	err := repo.DB.First(item, itemID).Error
	if err != nil {
		tx.Rollback()
		return nil
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (repo *Repository) Update(order *model.Order) error {
	tx := repo.DB.Begin()
	result := repo.DB.Model(order).Where("order_id = ?", order.OrderID).Updates(order)

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("no rows updated, check if the ID exists")
	}

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (repo *Repository) PlaceOrder(order *model.Order, CombineOrderItem *model.CombineOrderItem) error {
	tx := repo.DB.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error creating order: %v", err)
	}

	for _, orderedItem := range CombineOrderItem.Items {
		var existingItem model.Item
		if err := tx.Where("item_id = ?", orderedItem.ItemId).First(&existingItem).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				newItem := model.Item{
					ItemId: orderedItem.ItemId,
				}

				if err := tx.Create(&newItem).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("error creating item with ID %d: %v", orderedItem.ItemId, err)
				}
				existingItem = newItem
			} else {
				tx.Rollback()
				return fmt.Errorf("error checking item existence: %v", err)
			}
		}

		orderItem := model.OrderItem{
			OrderID:  order.OrderID,
			ItemId:   existingItem.ItemId,
			Quantity: orderedItem.Quantity,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error adding order item: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (repo *Repository) FetchAllOrder(orders *[]model.Order) error {
	tx := repo.DB.Begin()
	err := repo.DB.Find(orders).Error
	if err != nil {
		tx.Rollback()
		return nil
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (repo *Repository) FetchAverageOrderValueOfUser(input uint) (model.UserAverageOrderValue, error) {
	var result model.UserAverageOrderValue

	err := repo.DB.Table("orders").
		Select("user_id, AVG(total_bill) as average_order_value").
		Where("user_id = ?", input).
		Group("user_id").
		Scan(&result).Error

	if err != nil {
		return model.UserAverageOrderValue{}, err
	}
	return result, nil
}

func (repo *Repository) FetchAverageOrderValueOfRestaurant(input uint) (model.RestaurantAverageOrderValue, error) {
	var result model.RestaurantAverageOrderValue
	err := repo.DB.Table("orders").
		Select("restaurant_id, Avg(total_bill) as average_order_value").
		Where("restaurant_id = ?", input).
		Group("restaurant_id").
		Scan(&result).Error
	if err != nil {
		return model.RestaurantAverageOrderValue{}, err
	}
	return result, nil
}

func (repo *Repository) FetchAverageOrderValueBetweenTime(startTime string, endTime string) (model.TimeAverageOrderValue, error) {
	var result model.TimeAverageOrderValue
	err := repo.DB.Table("orders").
		Select("time, Avg(total_bill) as average_order_value").
		Where("time between ? and ?", startTime, endTime).
		Group("time").
		Scan(&result).Error
	if err != nil {
		return model.TimeAverageOrderValue{}, err
	}
	return result, nil
}

func (repo *Repository) FetchCompletedDeliversOfRider() ([]model.CompletedDelivers, error) {
	var result []model.CompletedDelivers
	err := repo.DB.Table("orders").
		Select("delivery_driver, count(*) as completed_delivers").
		Where("order_status = 'Completed' AND delivery_driver != 0").
		Group("delivery_driver").
		Order("completed_delivers DESC").
		Scan(&result).Error
	if err != nil {
		return []model.CompletedDelivers{}, err
	}
	log.Print(result[0].DeliveryDriver)
	return result, nil
}

func (repo *Repository) FetchCancelledOrdersWithItemDetails(Limit int, offset int) ([]model.OrderDetails, error) {
	var result []model.OrderDetails

	err := repo.DB.Table("orders").
		Select("orders.order_id, orders.user_id, orders.restaurant_id, STRING_AGG(order_items.item_id || ':' || order_items.quantity, ', ') AS item_details, orders.total_bill, orders.delivery_driver, orders.order_status, orders.time AS Order_time").
		Joins("inner join order_items on orders.order_id = order_items.order_id").
		Where("orders.order_status = ?", "Cancelled").
		Group("orders.order_id, orders.user_id, orders.restaurant_id, orders.total_bill, orders.delivery_driver, orders.order_status, orders.time").
		Limit(Limit).
		Offset(offset).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *Repository) FetchUserOrdersWithItemDetails(userID, limit, offset int) ([]model.OrderDetails, error) {
	var result []model.OrderDetails

	err := repo.DB.Table("orders").
		Select("orders.order_id, orders.user_id, orders.restaurant_id, STRING_AGG(order_items.item_id || ':' || order_items.quantity, ', ') AS item_details, orders.total_bill, orders.delivery_driver, orders.order_status, orders.time AS Order_time").
		Joins("inner join order_items on orders.order_id = order_items.order_id").
		Where("orders.user_id = ?", 1).
		Group("orders.order_id, orders.user_id, orders.restaurant_id, orders.total_bill, orders.delivery_driver, orders.order_status, orders.time").
		Limit(limit).
		Offset(offset).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *Repository) FetchTopPurchasedItems() ([]model.MostPurchasedItem, error) {
	var result []model.MostPurchasedItem
	err := repo.DB.Table("order_items").
		Select("orders.restaurant_id,order_items.item_id, COUNT(*) AS purchase_count ").
		Joins("INNER JOIN orders ON orders.order_id = order_items.order_id").
		Group("order_items.item_id, orders.restaurant_id").
		Order("purchase_count DESC").
		Limit(5).
		Scan(&result).Error

	if err != nil {
		return []model.MostPurchasedItem{}, err

	}
	return result, nil
}

func (repo *Repository) FetchCompletedOrdersCountByRestaurant(timeRange model.TimeRange) ([]model.RestaurantCompletedOrdersCount, error) {
	var result []model.RestaurantCompletedOrdersCount
	err := repo.DB.Table("orders").
		Select("restaurant_id, COUNT(*) AS completed_orders").
		Where("order_status = ? AND time BETWEEN ? AND ?", "Completed", timeRange.StartTime, timeRange.EndTime).
		Group("restaurant_id").
		Scan(&result).Error

	if err != nil {
		return []model.RestaurantCompletedOrdersCount{}, err
	}

	return result, nil

}

func (repo *Repository) FetchOrderStatusFrequencies() ([]model.OrderStatusFrequency, error) {
	var result []model.OrderStatusFrequency
	err := repo.DB.Table("orders").
		Select("order_status, COUNT(*) AS status_frequency").
		Group("order_status").
		Order("status_frequency DESC").
		Scan(&result).Error

	if err != nil {
		return []model.OrderStatusFrequency{}, err
	}

	return result, nil
}

func (repo *Repository) FetchTopFiveCustomers() ([]model.UserOrderFrequency, error) {
	var result []model.UserOrderFrequency
	err := repo.DB.Table("orders").
		Select("user_id, COUNT(*) AS order_frequency").
		Group("user_id").
		Order("order_frequency DESC").
		Limit(5).
		Scan(&result).Error

	if err != nil {
		return []model.UserOrderFrequency{}, err
	}

	return result, nil
}

func (repo *Repository) FetchRestaurantsRevenue() ([]model.RestaurantRevenue, error) {
	var result []model.RestaurantRevenue
	err := repo.DB.Table("orders").
		Select("restaurant_id , sum(total_bill) as revenue").
		Group("restaurant_id").
		Order("revenue DESC").
		Scan(&result).Error

	if err != nil {
		return []model.RestaurantRevenue{}, err
	}

	return result, nil
}

func (repo *Repository) FetchOrdersByDay(timeRange model.TimeFrame) ([]model.OrdersByDay, error) {
	var results []model.OrdersByDay
	err := repo.DB.Table("orders").
		Select("TO_CHAR(time, 'HH12:MI') AS hours_of_date,COUNT(total_bill) AS total_orders").
		Where("time BETWEEN ? AND ?", timeRange.StartDate, timeRange.EndDate).
		Group("hours_of_date").
		Order("total_orders DESC").
		Scan(&results).Error

	if err != nil {
		return []model.OrdersByDay{}, err
	}

	return results, nil
}
func (repo *Repository) FetchOrdersByWeek(timeRange model.TimeFrame) ([]model.OrdersByWeek, error) {
	var results []model.OrdersByWeek
	err := repo.DB.Table("orders").
		Select("TO_CHAR(time, 'DD') AS date_of_week,COUNT(total_bill) AS total_orders").
		Where("time BETWEEN ? AND ?", timeRange.StartDate, timeRange.EndDate).
		Group("date_of_week").
		Order("total_orders DESC").
		Scan(&results).Error

	if err != nil {
		return []model.OrdersByWeek{}, err
	}

	return results, nil
}

func (repo *Repository) FetchOrdersByMonth(timeRange model.TimeFrame) ([]model.OrderByMonth, error) {
	var results []model.OrderByMonth
	err := repo.DB.Table("orders").
		Select("EXTRACT(WEEK FROM time) - EXTRACT(WEEK FROM DATE_TRUNC('month', time)) + 1 AS week_of_month,COUNT(total_bill) AS total_orders").
		Where("time BETWEEN ? AND ?", timeRange.StartDate, timeRange.EndDate).
		Group("week_of_month").
		Order("total_orders DESC").
		Scan(&results).Error

	if err != nil {
		return []model.OrderByMonth{}, err
	}

	return results, nil
}

func (repo *Repository) FetchOrdersByYear(timeRange model.TimeFrame) ([]model.OrderByYear, error) {
	var results []model.OrderByYear
	err := repo.DB.Table("orders").
		Select("TO_CHAR(time, 'MM') AS month_of_year,COUNT(total_bill) AS total_orders").
		Where("time BETWEEN ? AND ?", timeRange.StartDate, timeRange.EndDate).
		Group("month_of_year").
		Order("total_orders DESC").
		Scan(&results).Error

	if err != nil {
		return []model.OrderByYear{}, err
	}

	return results, nil
}

func (repo *Repository) FetchOrderStatus(orderId uint) (model.ProcessOrder, error) {
	var result model.ProcessOrder
	err := repo.DB.Table("orders").
		Select("order_status").
		Where("order_id = ?", orderId).
		Scan(&result).Error

	if err != nil {
		return model.ProcessOrder{}, err
	}

	return result, nil
}

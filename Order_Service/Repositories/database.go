package database

import (
	"fmt"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	"gorm.io/gorm"
)

// Repository struct to handle dependency injection
type Repository struct {
	DB *gorm.DB
}

// NewRepository is a constructor function to initialize the repository with a DB connection
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (repo *Repository) GetOrders(order *[]model.Order, UserId int) error {
	err := repo.DB.Preload("Item").Where("user_id = ?", UserId).Find(order).Error
	return err
}

func (repo *Repository) GetOrder(order *model.Order, OrderId int) error {
	err := repo.DB.Where("order_id = ?", OrderId).First(order).Error
	return err
}

func (repo *Repository) Update(Model *model.Order, updateOrder payload.Order) error {
	// Generate dynamic update query using GORM
	result := repo.DB.Model(Model).Where("order_id = ?", updateOrder.OrderID).Updates(updateOrder)

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows updated, check if the ID exists")
	}

	// Check for errors
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *Repository) PutOrder(order *model.Order, CombineOrderItem *payload.CombineOrderItem) error {
	tx := repo.DB.Begin()

	order.UserId = CombineOrderItem.Order.UserId
	order.RestaurantID = CombineOrderItem.Order.RestaurantID
	order.TotalBill = 0
	order.OrderStatus = CombineOrderItem.Order.OrderStatus

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	totalBill := uint(0)

	for _, item := range CombineOrderItem.Items {

		if item.ItemId == 0 {
			tx.Rollback()
			return fmt.Errorf("Invalid ItemId: ItemId cannot be 0")
		}

		var existingItem model.Item
		if err := tx.Where("item_id = ?", item.ItemId).First(&existingItem).Error; err != nil {
			newItem := model.Item{
				ItemId:   item.ItemId,
				ItemName: item.ItemName,
			}
			if err := tx.Create(&newItem).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		itemTotal := item.ItemPrice * item.Quantity
		totalBill += itemTotal

		orderItem := model.OrderItem{
			OrderID:      order.OrderID,
			ItemId:       item.ItemId,
			RestaurantID: order.RestaurantID,
			ItemPrice:    item.ItemPrice,
			Quantity:     item.Quantity,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Model(&order).Update("TotalBill", totalBill).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

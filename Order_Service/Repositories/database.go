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

func (repo *Repository) GetOrders(order *[]model.Order, ID int, columnName string, orderDirection string) error {
	if orderDirection != "asc" && orderDirection != "desc" {
		orderDirection = "asc"
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

	err := repo.DB.Preload("Item").Where("user_id = ?", ID).Order(fmt.Sprintf("%s %s", columnName, orderDirection)).Find(order).Error
	return err
}

func (repo *Repository) GetOrder(order *model.Order, OrderId int) error {
	err := repo.DB.Where("order_id = ?", OrderId).First(order).Error
	return err
}

func (repo *Repository) Update(Model *model.Order, updateOrder payload.Order) error {
	result := repo.DB.Model(Model).Where("order_id = ?", updateOrder.OrderID).Updates(updateOrder)

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows updated, check if the ID exists")
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *Repository) PlaceOrder(order *model.Order, CombineOrderItem *payload.CombineOrderItem) error {
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

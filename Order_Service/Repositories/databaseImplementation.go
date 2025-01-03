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

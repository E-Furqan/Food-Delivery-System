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
	err := repo.DB.Where("UserId = ?", UserId).Find(order).Error
	return err
}

func (repo *Repository) Update(Model *model.Order, updateOrder payload.Order) error {
	// Generate dynamic update query using GORM
	result := repo.DB.Model(Model).Where("OrderID = ?", updateOrder.OrderID).Updates(updateOrder)

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

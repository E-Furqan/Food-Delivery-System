package database

import (
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
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

// CreateUser inserts a new user into the database
func (repo *Repository) CreateRestaurant(Restaurant *model.Restaurant) error {
	result := repo.DB.Create(Restaurant)
	log.Print(Restaurant)
	return result.Error
}

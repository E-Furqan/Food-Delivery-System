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
	InsertSourceConfiguration(configuration *model.CombinedStorageConfig) error
	InsertDestinationConfiguration(configuration *model.CombinedStorageConfig) error
}

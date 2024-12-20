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
	InsertConfiguration(configuration *model.Configs) error
	InsertToken(token *model.Token, configID uint) error
	LinkTokenToConfig(tokenID uint, configID uint) error
	DeleteExpiredTokens()
}

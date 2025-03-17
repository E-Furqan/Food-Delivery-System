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
	CreateRestaurant(restaurant *model.Restaurant) error
	GetRestaurant(columnName string, findParameter interface{}, restaurant *model.Restaurant) error
	LoadRestaurantWithItems(restaurant *model.Restaurant) error
	LoadItems(restaurantID uint, columnName string, order string) ([]model.Item, error)
	AddItemToRestaurantMenu(newItem model.Item) error
	RemoveItem(restaurantId uint, itemId uint) error
	GetAllRestaurants(restaurants *[]model.Restaurant) error
	UpdateRestaurantStatus(restaurant *model.Restaurant, input model.Input) error
	FetchOpenRestaurant() (model.OpenRestaurantCount, error)
	FetchItemPrices(items model.CombinedItemsRestaurantID) ([]model.Item, error)
}

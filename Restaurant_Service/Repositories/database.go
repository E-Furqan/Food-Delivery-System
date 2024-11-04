package database

import (
	"fmt"
	"log"

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

func (repo *Repository) CreateRestaurant(Restaurant *model.Restaurant) error {

	tx := repo.DB.Begin()

	result := tx.Create(Restaurant)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	log.Print(Restaurant)
	return tx.Commit().Error
}

func (repo *Repository) GetRestaurant(columnName string, findParameter interface{}, Restaurant *model.Restaurant) error {
	tx := repo.DB.Begin()

	err := repo.DB.Where(columnName+" = ?", findParameter).First(Restaurant).Error
	if err != nil {
		log.Printf("Error : %s", err)
		tx.Rollback()
		return err
	}

	err = repo.LoadRestaurantWithItems(Restaurant)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (repo *Repository) LoadRestaurantWithItems(Restaurant *model.Restaurant) error {

	tx := repo.DB.Begin()

	err := tx.Preload("Items").First(Restaurant, Restaurant.RestaurantId).Error
	if err != nil {
		log.Printf("Error loading Restaurant with Item: %v", err)
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *Repository) LoadItems(RestaurantID uint, columnName string, order string) ([]model.Item, error) {
	if columnName == "" {
		columnName = "restaurant_id"
	}
	if order == "" {
		order = "asc"
	}

	var ItemData []model.Item
	tx := repo.DB.Begin()

	err := tx.
		Where("restaurant_id = ?", RestaurantID).
		Order(fmt.Sprintf("%s %s", columnName, order)).
		Find(&ItemData).Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return ItemData, tx.Commit().Error
}

func (repo *Repository) AddItemToRestaurantMenu(newItem model.Item) error {

	tx := repo.DB.Begin()

	if err := tx.Create(&newItem).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (repo *Repository) RemoveItem(restaurantId uint, itemId uint) error {
	tx := repo.DB.Begin()

	err := tx.Where("restaurant_id = ? AND item_id = ?", restaurantId, itemId).Delete(&model.Item{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where("item_id = ?", itemId).Delete(&model.Item{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (repo *Repository) GetAllRestaurants(restaurants *[]model.Restaurant) error {
	err := repo.DB.
		Where("restaurant_status != ?", "closed").
		Find(restaurants).Error

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) UpdateRestaurantStatus(restaurant *model.Restaurant, input model.Input) error {

	tx := repo.DB.Begin()

	if err := tx.Model(restaurant).Update("restaurant_status", input.RestaurantStatus).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := repo.LoadRestaurantWithItems(restaurant)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

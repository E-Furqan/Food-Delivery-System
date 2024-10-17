package database

import (
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
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

func (repo *Repository) GetRestaurant(columnName string, findParameter interface{}, user *model.Restaurant) error {
	tx := repo.DB.Begin()

	err := repo.DB.Where(columnName+" = ?", findParameter).First(user).Error
	if err != nil {
		log.Printf("Error : %s", err)
		tx.Rollback()
		return err
	}

	err = repo.LoadRestaurantWithItems(user)
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

	log.Printf("Successfully loaded user with Item: %v", Restaurant.Items)
	return tx.Commit().Error
}

func (repo *Repository) LoadItemsInOrder(RestaurantID uint, columnName string, order string) ([]model.Item, error) {
	if columnName == "" {
		columnName = "user_id"
	}
	if order == "" {
		order = "asc"
	}

	var ItemData []model.Item
	tx := repo.DB.Begin()

	err := tx.Where("RestaurantID = ?", RestaurantID).
		Order(columnName + " " + order).Find(&ItemData).Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return ItemData, tx.Commit().Error
}

func (repo *Repository) AddItemToRestaurantMenu(restaurantId uint, newItem model.Item) error {

	tx := repo.DB.Begin()

	if err := tx.Create(&newItem).Error; err != nil {
		tx.Rollback()
		return err
	}

	restaurantItem := model.RestaurantItem{
		RestaurantId: restaurantId,
		ItemId:       newItem.ItemId,
	}

	if err := tx.Create(&restaurantItem).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *Repository) RemoveItemFromRestaurantMenu(restaurantId uint, itemId uint) error {
	tx := repo.DB.Begin()

	err := tx.Where("restaurant_id = ? AND item_id = ?", restaurantId, itemId).Delete(&model.RestaurantItem{}).Error
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

func (repo *Repository) GetAllRestaurants(restaurant *[]model.Restaurant) error {
	tx := repo.DB.Begin()

	err := tx.Where("restaurant_status != ?", "closed").Find(restaurant).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *Repository) UpdateRestaurantStatus(restaurant *model.Restaurant, input payload.Input) error {

	tx := repo.DB.Begin()

	if err := tx.Model(restaurant).Update("RestaurantStatus", input.RestaurantStatus).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

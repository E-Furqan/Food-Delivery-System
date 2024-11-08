package database_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println("Setup failed:", err)
		os.Exit(1)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	DB = TestDatabaseConnection()

	return m.Run(), nil
}

func setupRepository() *database.Repository {
	return &database.Repository{DB: DB}
}

func cleanupDB() {
	DB.Exec("DELETE FROM restaurants")
	DB.Exec("DELETE FROM items")
}

func TestCreateRestaurant(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	mockRestaurant := model.Restaurant{
		RestaurantName:        "RestaurantName",
		RestaurantAddress:     "Xyz Rawalpindi",
		RestaurantPhoneNumber: "123456789",
		RestaurantEmail:       "restaurant@gmail.com",
		Password:              "password123",
		RestaurantStatus:      "open",
	}

	err := repo.CreateRestaurant(&mockRestaurant)
	assert.NoError(t, err)

	var restaurant model.Restaurant

	err = DB.Where("restaurant_id", mockRestaurant.RestaurantId).Find(&restaurant).Error
	assert.NoError(t, err)

	assert.Equal(t, mockRestaurant.RestaurantName, restaurant.RestaurantName)
	assert.Equal(t, mockRestaurant.RestaurantAddress, restaurant.RestaurantAddress)
	assert.Equal(t, mockRestaurant.RestaurantPhoneNumber, restaurant.RestaurantPhoneNumber)
	assert.Equal(t, mockRestaurant.RestaurantEmail, restaurant.RestaurantEmail)
	assert.Equal(t, mockRestaurant.Password, restaurant.Password)
	assert.Equal(t, mockRestaurant.RestaurantStatus, restaurant.RestaurantStatus)

	err = repo.CreateRestaurant(&mockRestaurant)
	assert.Error(t, err)

	cleanupDB()
}

func TestGetRestaurant(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	mockRestaurant := model.Restaurant{
		RestaurantName:        "RestaurantName",
		RestaurantAddress:     "Xyz Rawalpindi",
		RestaurantPhoneNumber: "123456789",
		RestaurantEmail:       "restaurant@gmail.com",
		Password:              "password123",
		RestaurantStatus:      "open",
	}

	err := DB.Create(&mockRestaurant).Error
	assert.NoError(t, err)

	var restaurant model.Restaurant

	err = repo.GetRestaurant("restaurant_id", mockRestaurant.RestaurantId, &restaurant)
	assert.NoError(t, err)

	assert.Equal(t, mockRestaurant.RestaurantName, restaurant.RestaurantName)
	assert.Equal(t, mockRestaurant.RestaurantAddress, restaurant.RestaurantAddress)
	assert.Equal(t, mockRestaurant.RestaurantPhoneNumber, restaurant.RestaurantPhoneNumber)
	assert.Equal(t, mockRestaurant.RestaurantEmail, restaurant.RestaurantEmail)
	assert.Equal(t, mockRestaurant.Password, restaurant.Password)
	assert.Equal(t, mockRestaurant.RestaurantStatus, restaurant.RestaurantStatus)

	err = repo.GetRestaurant("restaurant_id", 96, &restaurant)
	assert.Error(t, err)

	cleanupDB()
}

func TestLoadRestaurantWithItems(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	mockRestaurant := model.Restaurant{
		RestaurantName:        "RestaurantName",
		RestaurantAddress:     "Xyz Rawalpindi",
		RestaurantPhoneNumber: "123456789",
		RestaurantEmail:       "restaurant@gmail.com",
		Password:              "password123",
		RestaurantStatus:      "open",
	}

	err := DB.Create(&mockRestaurant).Error
	assert.NoError(t, err)

	mockItems := []model.Item{
		{ItemName: "Pizza",
			ItemDescription: "Pizza with stuffed crust",
			ItemPrice:       150.0,
			RestaurantId:    mockRestaurant.RestaurantId},
		{ItemName: "burger",
			ItemDescription: "grilled beef burger",
			ItemPrice:       200.0,
			RestaurantId:    mockRestaurant.RestaurantId},
	}

	for _, mockItem := range mockItems {
		err := DB.Create(&mockItem).Error
		assert.NoError(t, err)
	}

	err = repo.LoadRestaurantWithItems(&mockRestaurant)
	assert.NoError(t, err)

	assert.NotEmpty(t, mockRestaurant.Items)
	assert.Equal(t, len(mockItems), len(mockRestaurant.Items))

	for i, item := range mockRestaurant.Items {
		assert.Equal(t, mockItems[i].ItemName, item.ItemName)
		assert.Equal(t, mockItems[i].ItemDescription, item.ItemDescription)
		assert.Equal(t, mockItems[i].ItemPrice, item.ItemPrice)
		assert.Equal(t, mockItems[i].RestaurantId, item.RestaurantId)
	}

	cleanupDB()
}

func TestLoadItems(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	mockRestaurant := model.Restaurant{
		RestaurantName:        "RestaurantName",
		RestaurantAddress:     "Xyz Rawalpindi",
		RestaurantPhoneNumber: "123456789",
		RestaurantEmail:       "restaurant@gmail.com",
		Password:              "password123",
		RestaurantStatus:      "open",
	}

	err := DB.Create(&mockRestaurant).Error
	assert.NoError(t, err)

	mockItems := []model.Item{
		{
			ItemName:        "Pizza",
			ItemDescription: "Pizza with stuffed crust",
			ItemPrice:       150.0,
			RestaurantId:    mockRestaurant.RestaurantId,
		},
		{
			ItemName:        "Burger",
			ItemDescription: "Grilled beef burger",
			ItemPrice:       200.0,
			RestaurantId:    mockRestaurant.RestaurantId,
		},
	}

	for _, mockItem := range mockItems {
		err := DB.Create(&mockItem).Error
		assert.NoError(t, err)
	}

	var items []model.Item
	items, err = repo.LoadItems(mockRestaurant.RestaurantId, "", "")
	assert.NoError(t, err)

	assert.NotEmpty(t, items)
	assert.Equal(t, len(mockItems), len(items))

	for i, item := range items {
		assert.Equal(t, mockItems[i].ItemName, item.ItemName)
		assert.Equal(t, mockItems[i].ItemDescription, item.ItemDescription)
		assert.Equal(t, mockItems[i].ItemPrice, item.ItemPrice)
		assert.Equal(t, mockItems[i].RestaurantId, item.RestaurantId)
	}

	cleanupDB()
}

func TestAddItemToRestaurantMenu(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	mockRestaurant := model.Restaurant{
		RestaurantName:        "Test Restaurant",
		RestaurantAddress:     "123 Test St",
		RestaurantPhoneNumber: "1234567890",
		RestaurantEmail:       "testrestaurant@example.com",
		Password:              "password123",
		RestaurantStatus:      "open",
	}

	err := DB.Create(&mockRestaurant).Error
	assert.NoError(t, err)

	newItem := model.Item{
		ItemName:        "Pasta",
		ItemDescription: "Delicious Italian Pasta",
		ItemPrice:       12.99,
		RestaurantId:    mockRestaurant.RestaurantId,
	}

	err = repo.AddItemToRestaurantMenu(newItem)
	assert.NoError(t, err)

	var addedItem model.Item
	err = DB.Where("item_name = ?", newItem.ItemName).First(&addedItem).Error
	assert.NoError(t, err)
	assert.Equal(t, newItem.ItemName, addedItem.ItemName)
	assert.Equal(t, newItem.ItemDescription, addedItem.ItemDescription)
	assert.Equal(t, newItem.ItemPrice, addedItem.ItemPrice)
	assert.Equal(t, newItem.RestaurantId, addedItem.RestaurantId)

	cleanupDB()
}

func TestRemoveItem(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	mockRestaurant := model.Restaurant{
		RestaurantName:        "Test Restaurant",
		RestaurantAddress:     "123 Test St",
		RestaurantPhoneNumber: "1234567890",
		RestaurantEmail:       "testrestaurant@example.com",
		Password:              "password123",
		RestaurantStatus:      "open",
	}

	err := DB.Create(&mockRestaurant).Error
	assert.NoError(t, err)

	mockItem := model.Item{
		ItemName:        "Pizza",
		ItemDescription: "Delicious cheese pizza",
		ItemPrice:       10.99,
		RestaurantId:    mockRestaurant.RestaurantId,
	}

	err = DB.Create(&mockItem).Error
	assert.NoError(t, err)

	err = repo.RemoveItem(mockRestaurant.RestaurantId, mockItem.ItemId)
	assert.NoError(t, err)

	var removedItem model.Item
	err = DB.Where("item_id = ?", mockItem.ItemId).First(&removedItem).Error
	assert.Error(t, err)

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected gorm.ErrRecordNotFound, got: %v", err)
	}

	cleanupDB()
}

func TestGetAllRestaurants(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	mockRestaurants := []model.Restaurant{
		{
			RestaurantName:        "Open Restaurant 1",
			RestaurantAddress:     "Address 1",
			RestaurantPhoneNumber: "1234567890",
			RestaurantEmail:       "open1@example.com",
			Password:              "password123",
			RestaurantStatus:      "open",
		},
		{
			RestaurantName:        "Closed Restaurant",
			RestaurantAddress:     "Address 2",
			RestaurantPhoneNumber: "0987654321",
			RestaurantEmail:       "closed@example.com",
			Password:              "password123",
			RestaurantStatus:      "closed",
		},
		{
			RestaurantName:        "Open Restaurant 2",
			RestaurantAddress:     "Address 3",
			RestaurantPhoneNumber: "5555555555",
			RestaurantEmail:       "open2@example.com",
			Password:              "password123",
			RestaurantStatus:      "open",
		},
	}

	for _, restaurant := range mockRestaurants {
		err := DB.Create(&restaurant).Error
		assert.NoError(t, err)
	}

	var retrievedRestaurants []model.Restaurant
	err := repo.GetAllRestaurants(&retrievedRestaurants)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(retrievedRestaurants))

	expectedNames := map[string]struct{}{
		"Open Restaurant 1": {},
		"Open Restaurant 2": {},
	}

	for _, restaurant := range retrievedRestaurants {
		_, exists := expectedNames[restaurant.RestaurantName]
		assert.True(t, exists, "Unexpected restaurant found: %s", restaurant.RestaurantName)
	}

	cleanupDB()
}

func TestUpdateRestaurantStatus(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	mockRestaurant := model.Restaurant{
		RestaurantName:        "Test Restaurant",
		RestaurantAddress:     "123 Test St",
		RestaurantPhoneNumber: "1234567890",
		RestaurantEmail:       "testrestaurant@example.com",
		Password:              "password123",
		RestaurantStatus:      "open",
	}

	err := DB.Create(&mockRestaurant).Error
	assert.NoError(t, err)

	input := model.Input{
		RestaurantStatus: "closed",
	}

	err = repo.UpdateRestaurantStatus(&mockRestaurant, input)
	assert.NoError(t, err)

	var updatedRestaurant model.Restaurant
	err = DB.First(&updatedRestaurant, mockRestaurant.RestaurantId).Error
	assert.NoError(t, err)
	assert.Equal(t, "closed", updatedRestaurant.RestaurantStatus)

	cleanupDB()
}

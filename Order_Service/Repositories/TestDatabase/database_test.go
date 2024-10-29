package database_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/stretchr/testify/assert"
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

func TestGetOrders(t *testing.T) {
	repo := setupRepository()

	mockOrder := model.Order{
		UserId:           1,
		RestaurantID:     1,
		DeliveryDriverID: 0,
		OrderStatus:      "completed",
		TotalBill:        50.0,
		Time:             time.Now(),
	}
	err := DB.Create(&mockOrder).Error
	assert.NoError(t, err)

	mockItem := model.Item{
		ItemId: 1,
	}
	err = DB.Create(&mockItem).Error
	assert.NoError(t, err)

	mockOrderItem := model.OrderItem{
		OrderID:  mockOrder.OrderID,
		ItemId:   mockItem.ItemId,
		Quantity: 2,
	}
	err = DB.Create(&mockOrderItem).Error
	assert.NoError(t, err)

	var orders []model.Order

	err = repo.GetOrders(&orders, mockOrder.UserId, "order_id", "asc", "user_id")
	assert.NoError(t, err)

	assert.Len(t, orders, 1)
	assert.Equal(t, mockOrder.OrderID, orders[0].OrderID)
	assert.Equal(t, mockOrder.OrderStatus, orders[0].OrderStatus)
	assert.Equal(t, mockOrder.TotalBill, orders[0].TotalBill)
	assert.Len(t, orders[0].Item, 1)
	assert.Equal(t, mockItem.ItemId, orders[0].Item[0].ItemId)

	cleanupDB()
}

func cleanupDB() {
	DB.Exec("DELETE FROM order_items")
	DB.Exec("DELETE FROM items")
	DB.Exec("DELETE FROM orders")
}

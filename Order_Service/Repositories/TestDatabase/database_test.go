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

func cleanupDB() {
	DB.Exec("DELETE FROM order_items")
	DB.Exec("DELETE FROM items")
	DB.Exec("DELETE FROM orders")
}

func TestGetOrders(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	mockOrder := model.Order{
		UserId:         1,
		RestaurantID:   1,
		DeliveryDriver: 0,
		OrderStatus:    "completed",
		TotalBill:      50.0,
		Time:           time.Now(),
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

	err = repo.GetOrders(&orders, mockOrder.UserId, "order_id", "invalid_sort_order", "user_id")
	assert.NoError(t, err)

	err = repo.GetOrders(&orders, mockOrder.UserId, "invalid_column", "asc", "user_id")
	assert.NoError(t, err)

	err = repo.GetOrders(&orders, 9999, "order_id", "asc", "user_id")
	assert.NoError(t, err)
	assert.Len(t, orders, 0)

	err = repo.GetOrders(&orders, mockOrder.UserId, "order_id", "asc", "invalid_search_column")
	assert.Error(t, err)

	cleanupDB()
}

func TestGetOrder(t *testing.T) {
	repo := setupRepository()
	cleanupDB()
	mockOrder := model.Order{
		UserId:         1,
		RestaurantID:   1,
		DeliveryDriver: 0,
		OrderStatus:    "completed",
		TotalBill:      50.0,
		Time:           time.Now(),
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

	var order model.Order

	err = repo.GetOrder(&order, mockOrder.OrderID)
	assert.NoError(t, err)

	assert.Equal(t, mockOrder.OrderID, order.OrderID)
	assert.Equal(t, mockOrder.OrderStatus, order.OrderStatus)
	assert.Equal(t, mockOrder.TotalBill, order.TotalBill)
	assert.Len(t, order.Item, 1)
	assert.Equal(t, mockItem.ItemId, order.Item[0].ItemId)
	cleanupDB()
	var order1 model.Order
	err = repo.GetOrder(&order1, 8989898)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), order1.OrderID)

	cleanupDB()
}

func TestGetOrderWithoutRider(t *testing.T) {
	repo := setupRepository()
	cleanupDB()
	mockOrder := model.Order{
		UserId:         1,
		RestaurantID:   1,
		DeliveryDriver: 0,
		OrderStatus:    "completed",
		TotalBill:      50.0,
		Time:           time.Now(),
	}
	err := DB.Create(&mockOrder).Error
	assert.NoError(t, err)

	mockOrder1 := model.Order{
		UserId:         2,
		RestaurantID:   1,
		DeliveryDriver: 1,
		OrderStatus:    "order placed",
		TotalBill:      500.0,
		Time:           time.Now(),
	}
	err = DB.Create(&mockOrder1).Error
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

	mockOrderItem1 := model.OrderItem{
		OrderID:  mockOrder1.OrderID,
		ItemId:   mockItem.ItemId,
		Quantity: 2,
	}
	err = DB.Create(&mockOrderItem1).Error
	assert.NoError(t, err)

	var orders []model.Order

	err = repo.GetOrderWithoutRider(&orders)
	assert.NoError(t, err)

	assert.Len(t, orders, 1)
	assert.Equal(t, mockOrder.OrderID, orders[0].OrderID)
	assert.Equal(t, mockOrder.OrderStatus, orders[0].OrderStatus)
	assert.Equal(t, mockOrder.TotalBill, orders[0].TotalBill)
	assert.Len(t, orders[0].Item, 1)
	assert.Equal(t, mockItem.ItemId, orders[0].Item[0].ItemId)
	assert.Equal(t, mockOrder.DeliveryDriver, orders[0].DeliveryDriver)

	cleanupDB()
}

func TestGetOrderItems(t *testing.T) {
	repo := setupRepository()
	cleanupDB()
	mockOrder := model.Order{
		UserId:         1,
		RestaurantID:   1,
		DeliveryDriver: 0,
		OrderStatus:    "completed",
		TotalBill:      50.0,
		Time:           time.Now(),
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

	var OrderItem []model.OrderItem

	err = repo.GetOrderItems(&OrderItem, mockOrder.OrderID)
	assert.NoError(t, err)

	assert.Len(t, OrderItem, 1)
	assert.Equal(t, mockOrder.OrderID, OrderItem[0].OrderID)
	assert.Equal(t, mockItem.ItemId, OrderItem[0].ItemId)
	assert.Equal(t, mockOrderItem.Quantity, OrderItem[0].Quantity)

	cleanupDB()
}

func TestGetItemByID(t *testing.T) {
	repo := setupRepository()
	cleanupDB()
	mockItem := model.Item{
		ItemId: 1,
	}
	err := DB.Create(&mockItem).Error
	assert.NoError(t, err)

	mockItem1 := model.Item{
		ItemId: 2,
	}
	err = DB.Create(&mockItem1).Error
	assert.NoError(t, err)

	var item model.Item

	err = repo.GetItemByID(mockItem.ItemId, &item)
	assert.NoError(t, err)

	var item1 model.Item
	err = repo.GetItemByID(mockItem1.ItemId, &item1)
	assert.NoError(t, err)

	assert.Equal(t, mockItem.ItemId, item.ItemId)
	assert.Equal(t, mockItem1.ItemId, item1.ItemId)

	cleanupDB()
}

func TestUpdate(t *testing.T) {
	repo := setupRepository()
	cleanupDB()
	mockOrder := model.Order{
		UserId:         1,
		RestaurantID:   1,
		DeliveryDriver: 0,
		OrderStatus:    "completed",
		TotalBill:      50.0,
		Time:           time.Now(),
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

	mockOrder.UserId = 2
	mockOrder.RestaurantID = 2
	mockOrder.DeliveryDriver = 3
	mockOrder.OrderStatus = "order placed"
	mockOrder.TotalBill = 100
	mockOrder.Time = time.Now()

	err = repo.Update(&mockOrder)
	assert.NoError(t, err)

	var order model.Order
	err = DB.Where("order_id = ?", mockOrder.OrderID).Preload("Item").First(&order).Error
	assert.NoError(t, err)

	assert.Equal(t, mockOrder.OrderID, order.OrderID)
	assert.Equal(t, mockOrder.UserId, order.UserId)
	assert.Equal(t, mockOrder.RestaurantID, order.RestaurantID)
	assert.Equal(t, mockOrder.DeliveryDriver, order.DeliveryDriver)
	assert.Equal(t, mockOrder.OrderStatus, order.OrderStatus)
	assert.Equal(t, mockOrder.TotalBill, order.TotalBill)
	expectedTime := mockOrder.Time.Truncate(time.Second)
	actualTime := order.Time.Truncate(time.Second)
	assert.Equal(t, expectedTime, actualTime)

	cleanupDB()
}

func TestPlaceOrder(t *testing.T) {
	repo := setupRepository()
	cleanupDB()
	mockOrder := model.Order{
		UserId:         1,
		RestaurantID:   1,
		DeliveryDriver: 0,
		OrderStatus:    "order placed",
		TotalBill:      50.0,
		Time:           time.Now(),
	}

	mockCombineOrderItem := model.CombineOrderItem{
		ID: model.ID{
			OrderID:        mockOrder.OrderID,
			RestaurantId:   mockOrder.RestaurantID,
			UserId:         mockOrder.UserId,
			DeliveryDriver: mockOrder.DeliveryDriver,
		},
		Items: []model.OrderItemPayload{
			{ItemId: 1, Quantity: 2},
			{ItemId: 2, Quantity: 2},
		},
	}

	err := repo.PlaceOrder(&mockOrder, &mockCombineOrderItem)
	assert.NoError(t, err)

	var order model.Order
	err = DB.Where("order_id = ?", mockOrder.OrderID).First(&order).Error
	assert.NoError(t, err)

	assert.Equal(t, mockOrder.OrderID, order.OrderID)
	assert.Equal(t, mockOrder.UserId, order.UserId)
	assert.Equal(t, mockOrder.RestaurantID, order.RestaurantID)
	assert.Equal(t, mockOrder.OrderStatus, order.OrderStatus)
	assert.Equal(t, mockOrder.TotalBill, order.TotalBill)
	assert.Equal(t, mockOrder.DeliveryDriver, order.DeliveryDriver)
	expectedTime := mockOrder.Time.Truncate(time.Second)
	actualTime := order.Time.Truncate(time.Second)
	assert.Equal(t, expectedTime, actualTime)

	var orderItems []model.OrderItem
	err = DB.Where("order_id = ?", mockOrder.OrderID).Find(&orderItems).Error
	assert.NoError(t, err)
	assert.Len(t, orderItems, len(mockCombineOrderItem.Items))

	for i, item := range mockCombineOrderItem.Items {
		assert.Equal(t, item.ItemId, orderItems[i].ItemId)
		assert.Equal(t, item.Quantity, orderItems[i].Quantity)
	}

	cleanupDB()
}

func TestFetchAllOrders(t *testing.T) {
	repo := setupRepository()
	cleanupDB()
	mockOrders := []model.Order{
		{
			UserId:         1,
			RestaurantID:   1,
			DeliveryDriver: 0,
			OrderStatus:    "completed",
			TotalBill:      50.0,
			Time:           time.Now(),
		},
		{
			UserId:         2,
			RestaurantID:   1,
			DeliveryDriver: 1,
			OrderStatus:    "pending",
			TotalBill:      75.0,
			Time:           time.Now(),
		},
	}
	for _, mockOrder := range mockOrders {
		err := DB.Create(&mockOrder).Error
		assert.NoError(t, err)
	}

	var fetchedOrders []model.Order
	err := repo.FetchAllOrder(&fetchedOrders)
	assert.NoError(t, err)

	assert.Len(t, fetchedOrders, len(mockOrders))
	for i, order := range fetchedOrders {
		assert.Equal(t, mockOrders[i].OrderStatus, order.OrderStatus)
		assert.Equal(t, mockOrders[i].TotalBill, order.TotalBill)
		assert.Equal(t, mockOrders[i].UserId, order.UserId)
		assert.Equal(t, mockOrders[i].RestaurantID, order.RestaurantID)
	}

	cleanupDB()
}

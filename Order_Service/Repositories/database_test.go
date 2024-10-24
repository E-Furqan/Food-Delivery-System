package database_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

// Test for GetOrders
func TestRepository_GetOrders(t *testing.T) {
	db, mock, err := setupMockDB()
	if err != nil {
		t.Fatalf("Failed to setup mock DB: %v", err)
	}

	repo := database.NewRepository(db)

	var orders []model.Order
	ID := 1
	columnName := "order_id"
	orderDirection := "asc"
	searchColumn := "user_id"

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT \\* FROM \"orders\" WHERE \"user_id\" = \\$1").
		WithArgs(ID).
		WillReturnRows(sqlmock.NewRows([]string{"order_id", "user_id", "restaurant_id", "order_status"}).
			AddRow(1, 1, 1, "pending"))
	mock.ExpectCommit()

	err = repo.GetOrders(&orders, ID, columnName, orderDirection, searchColumn)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(orders))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

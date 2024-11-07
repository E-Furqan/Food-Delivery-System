package config

import (
	"fmt"
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	model.DatabaseConfigEnv
}

func NewDatabase(env model.DatabaseConfigEnv) *DatabaseConfig {
	return &DatabaseConfig{
		DatabaseConfigEnv: env,
	}
}

var DB *gorm.DB

func (DatabaseConfig *DatabaseConfig) Connection() *gorm.DB {

	var connection_string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DatabaseConfig.DatabaseConfigEnv.DATABASE_HOST, DatabaseConfig.DatabaseConfigEnv.DATABASE_PORT,
		DatabaseConfig.DatabaseConfigEnv.DATABASE_USER, DatabaseConfig.DatabaseConfigEnv.DATABASE_PASSWORD,
		DatabaseConfig.DatabaseConfigEnv.DATABASE_NAME)

	DB, err := gorm.Open(postgres.Open(connection_string), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&model.Order{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = DB.AutoMigrate(&model.Item{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = DB.AutoMigrate(&model.OrderItem{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connection established")

	return DB
}

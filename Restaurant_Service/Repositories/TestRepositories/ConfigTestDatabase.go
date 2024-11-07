package database_test

import (
	"fmt"
	"log"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func TestDatabaseConnection() *gorm.DB {
	envVar := environmentVariable.ReadDatabaseEnv()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var connection_string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		envVar.DATABASE_HOST, envVar.DATABASE_PORT, envVar.DATABASE_USER, envVar.DATABASE_PASSWORD, "testrestaurant")

	DB, err = gorm.Open(postgres.Open(connection_string), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&model.Restaurant{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = DB.AutoMigrate(&model.Item{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database connection established")

	return DB
}

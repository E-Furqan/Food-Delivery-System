package database_test

import (
	"fmt"
	"log"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func TestDatabaseConnection() *gorm.DB {
	envVar := environmentVariable.ReadDatabaseEnv()

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	var connection_string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"0.0.0.0", 5432, envVar.DATABASE_USER, envVar.DATABASE_PASSWORD, "testrestaurant")

	DB, err = gorm.Open(postgres.Open(connection_string), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return DB
}

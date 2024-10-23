package config

import (
	"fmt"
	"log"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	user     string
	password string
	host     string
	db_name  string
	port     int
)

var DB *gorm.DB

func SetEnvValue(envVar environmentVariable.Environment) {
	user = envVar.USER
	password = envVar.PASSWORD
	host = envVar.HOST
	db_name = envVar.DB_NAME
	port = envVar.PORT
}

func Connection() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var connection_string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, db_name)

	DB, err = gorm.Open(postgres.Open(connection_string), &gorm.Config{})
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

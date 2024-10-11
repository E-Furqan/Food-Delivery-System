package config

import (
	"fmt"
	"log"
	"strconv"

	environmentvariable "github.com/E-Furqan/Food-Delivery-System/enviorment_variable"
	model "github.com/E-Furqan/Food-Delivery-System/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connection() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	portStr := environmentvariable.Get_env("PORT")
	port, erro := strconv.Atoi(portStr)

	if erro != nil {
		log.Fatalf("Failed to convert port to integer : %v", erro)
	}

	var (
		user     = environmentvariable.Get_env("USER1")
		password = environmentvariable.Get_env("PASSWORD")
		host     = environmentvariable.Get_env("HOST")
		db_name  = environmentvariable.Get_env("DB_NAME")
	)
	var connectionstring = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, db_name)

	DB, err = gorm.Open(postgres.Open(connectionstring), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&model.Role{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connection established")

	return DB
}

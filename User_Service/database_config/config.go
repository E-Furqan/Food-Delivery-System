package config

import (
	"fmt"
	"log"

	entity "github.com/E-Furqan/Food-Delivery-System/Entity"
	environmentvariable "github.com/E-Furqan/Food-Delivery-System/enviorment_variable"
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
	envVar := environmentvariable.ReadEnv()
	port := envVar.PORT

	var (
		user     = envVar.USER
		password = envVar.PASSWORD
		host     = envVar.HOST
		db_name  = envVar.DB_NAME
	)
	var connectionstring = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, db_name)

	DB, err = gorm.Open(postgres.Open(connectionstring), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&entity.Role{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = DB.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connection established")

	return DB
}

package config

import (
	"fmt"
	"log"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	environmentVariable.Environment
}

func NewDatabase(env environmentVariable.Environment) *DatabaseConfig {
	return &DatabaseConfig{
		Environment: env,
	}
}

var DB *gorm.DB

func (DatabaseConfig *DatabaseConfig) Connection() *gorm.DB {

	var connection_string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DatabaseConfig.Environment.HOST, DatabaseConfig.Environment.PORT, DatabaseConfig.Environment.USER, DatabaseConfig.Environment.PASSWORD, DatabaseConfig.Environment.DB_NAME)

	DB, err := gorm.Open(postgres.Open(connection_string), &gorm.Config{})
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

	err = DB.AutoMigrate(&model.UserRole{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connection established")

	return DB
}

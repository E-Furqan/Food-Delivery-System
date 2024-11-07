package config

import (
	"fmt"
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	model.DatabaseEnv
}

func NewDatabase(env model.DatabaseEnv) *DatabaseConfig {
	return &DatabaseConfig{
		DatabaseEnv: env,
	}
}

var DB *gorm.DB

func (DatabaseConfig *DatabaseConfig) Connection() *gorm.DB {

	var connection_string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DatabaseConfig.DatabaseEnv.DATABASE_HOST, DatabaseConfig.DatabaseEnv.DATABASE_PORT, DatabaseConfig.DatabaseEnv.DATABASE_USER, DatabaseConfig.DatabaseEnv.DATABASE_PASSWORD, DatabaseConfig.DatabaseEnv.DATABASE_NAME)

	DB, err := gorm.Open(postgres.Open(connection_string), &gorm.Config{})
	log.Print(connection_string)
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

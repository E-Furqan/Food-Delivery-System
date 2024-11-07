package config

import (
	"fmt"
	"log"
	"path/filepath"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var connection_string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DatabaseConfig.DatabaseEnv.DATABASE_HOST, DatabaseConfig.DatabaseEnv.DATABASE_PORT, DatabaseConfig.DatabaseEnv.DATABASE_USER,
		DatabaseConfig.DatabaseEnv.DATABASE_PASSWORD, DatabaseConfig.DatabaseEnv.DATABASE_NAME)

	log.Println(connection_string)
	DB, err = gorm.Open(postgres.Open(connection_string), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// err = DB.AutoMigrate(&model.Item{})
	// if err != nil {
	// 	log.Fatalf("Failed to migrate database: %v", err)
	// }

	// err = DB.AutoMigrate(&model.Restaurant{})
	// if err != nil {
	// 	log.Fatalf("Failed to migrate database: %v", err)
	// }

	log.Println("Database connection established")

	return DB
}

func (DatabaseConfig *DatabaseConfig) RunMigrations() {
	// Use Docker service name for the database connection
	connectionString := "postgres://furqan:furqan@restaurant_service_db_1:5432/Restaurant?sslmode=disable"

	absPath, err := filepath.Abs("./Migration")
	if err != nil {
		log.Printf("Error getting absolute path: %v", err)
	}
	log.Printf("pringint %s", absPath)
	absPath = "file://" + absPath
	log.Printf("Using migration path: %s", absPath)
	log.Printf("Using connectionString: %s", connectionString)

	m, err := migrate.New(
		absPath, // Ensure the source has file:// scheme
		connectionString,
	)
	if err != nil {
		log.Fatalf("Migration instance error: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations applied successfully")
}

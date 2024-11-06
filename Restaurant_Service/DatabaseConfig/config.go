package config

import (
	"fmt"
	"log"
	"path/filepath"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var connection_string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DatabaseConfig.Environment.DATABASE_HOST, DatabaseConfig.Environment.DATABASE_PORT, DatabaseConfig.Environment.DATABASE_USER,
		DatabaseConfig.Environment.DATABASE_PASSWORD, DatabaseConfig.Environment.DATABASE_NAME)

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
	// Use file:// scheme for local migrations
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		DatabaseConfig.Environment.DATABASE_USER, DatabaseConfig.Environment.DATABASE_PASSWORD,
		DatabaseConfig.Environment.DATABASE_HOST, DatabaseConfig.Environment.DATABASE_PORT,
		DatabaseConfig.Environment.DATABASE_NAME)

	// Get the absolute path of the migration folder
	absPath, err := filepath.Abs("./Migration")
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}
	log.Printf("migration %s", absPath)
	// Initialize migrate instance with file:// scheme and the absolute path
	m, err := migrate.New(
		"file://"+absPath, // Ensure the source has file:// scheme
		connectionString,
	)
	if err != nil {
		log.Fatalf("Migration instance error: %v", err)
	}

	// Run migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations applied successfully")
}

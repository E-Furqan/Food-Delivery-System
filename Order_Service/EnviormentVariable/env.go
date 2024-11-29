package environmentVariable

import (
	"log"
	"strconv"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/joho/godotenv"
)

func ReadDatabaseConfigEnv() model.DatabaseConfigEnv {
	var envVar model.DatabaseConfigEnv

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	envVar.DATABASE_HOST = utils.GetEnv("DATABASE_HOST", "localhost")
	envVar.DATABASE_USER = utils.GetEnv("DATABASE_USER", "furqan")
	envVar.DATABASE_PASSWORD = utils.GetEnv("DATABASE_PASSWORD", "furqan")
	envVar.DATABASE_NAME = utils.GetEnv("DATABASE_NAME", "Order")
	portStr := utils.GetEnv("DATABASE_PORT", "5431")
	envVar.DATABASE_PORT, err = strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Error converting PORT to integer: %v", err)
		envVar.DATABASE_PORT = 5432
	}
	return envVar
}

func ReadRestaurantClientEnv() model.RestaurantClientEnv {
	var envVar model.RestaurantClientEnv

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	envVar.BASE_URL = utils.GetEnv("BASE_URL", "http://localhost")
	envVar.Get_Items_URL = utils.GetEnv("Get_Items_URL", "/restaurant/view/menu")
	envVar.RESTAURANT_PORT = utils.GetEnv("RESTAURANT_PORT", ":8082")

	return envVar
}

func ReadMiddlewareEnv() model.MiddlewareEnv {
	var envVar model.MiddlewareEnv

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	envVar.JWT_SECRET = utils.GetEnv("JWT_SECRET", "Furqan")
	envVar.RefreshTokenKey = utils.GetEnv("REFRESH_TOKEN_SECRET", "Ali")
	return envVar
}

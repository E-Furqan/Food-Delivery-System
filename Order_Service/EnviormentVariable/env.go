package environmentVariable

import (
	"log"
	"strconv"

	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/joho/godotenv"
)

type Environment struct {
	DATABASE_PORT                int
	DATABASE_HOST                string
	DATABASE_USER                string
	DATABASE_PASSWORD            string
	DATABASE_NAME                string
	Get_Items_URL                string
	BASE_URL                     string
	Process_Order_Restaurant_URL string
	RESTAURANT_PORT              string
	USER_PORT                    string
	Process_Order_User_URL       string
	JWT_SECRET                   string
	RefreshTokenKey              string
}

// ReadEnv reads environment variables from a .env file and returns an Environment struct
func ReadEnv() Environment {
	var envVar Environment

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	envVar.DATABASE_HOST = utils.GetEnv("DATABASE_HOST", "db")
	envVar.DATABASE_USER = utils.GetEnv("DATABASE_USER", "furqan")
	envVar.DATABASE_PASSWORD = utils.GetEnv("DATABASE_PASSWORD", "furqan")
	envVar.DATABASE_NAME = utils.GetEnv("DATABASE_NAME", "Order")
	envVar.BASE_URL = utils.GetEnv("BASE_URL", "http://localhost")
	envVar.Get_Items_URL = utils.GetEnv("Get_Items_URL", "/restaurant/view/menu")
	envVar.Process_Order_Restaurant_URL = utils.GetEnv("Process_Order_Restaurant_URL", "/restaurant/process/order")
	envVar.Process_Order_User_URL = utils.GetEnv("Process_Order_User_URL", "/user/process/order")
	envVar.RESTAURANT_PORT = utils.GetEnv("RESTAURANT_PORT", ":8082")
	envVar.USER_PORT = utils.GetEnv("USER_PORT", ":8083")
	portStr := utils.GetEnv("DATABASE_PORT", "5432")
	envVar.DATABASE_PORT, err = strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Error converting PORT to integer: %v", err)
		envVar.DATABASE_PORT = 5432
	}
	envVar.JWT_SECRET = utils.GetEnv("JWT_SECRET", "Furqan")
	envVar.RefreshTokenKey = utils.GetEnv("REFRESH_TOKEN_SECRET", "Ali")
	return envVar
}

package environmentVariable

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment struct {
	PORT                         int
	HOST                         string
	USER                         string
	PASSWORD                     string
	DB_NAME                      string
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
		log.Fatalf("Error loading .env file: %v", err)
	}

	envVar.HOST = os.Getenv("HOST")
	envVar.USER = os.Getenv("USER1")
	envVar.PASSWORD = os.Getenv("PASSWORD")
	envVar.DB_NAME = os.Getenv("DB_NAME")
	envVar.BASE_URL = os.Getenv("BASE_URL")
	envVar.Get_Items_URL = os.Getenv("Get_Items_URL")
	envVar.Process_Order_Restaurant_URL = os.Getenv("Process_Order_Restaurant_URL")
	envVar.Process_Order_User_URL = os.Getenv("Process_Order_User_URL")
	envVar.RESTAURANT_PORT = os.Getenv("RESTAURANT_PORT")
	envVar.USER_PORT = os.Getenv("USER_PORT")
	portStr := os.Getenv("PORT")
	envVar.PORT, err = strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error converting PORT to integer: %v", err)
	}
	envVar.JWT_SECRET = os.Getenv("JWT_SECRET")
	envVar.RefreshTokenKey = os.Getenv("REFRESH_TOKEN_SECRET")
	return envVar
}

func Get_env(key string) string {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("environment variable %s is not set", key)
		return "" // Return an error if the variable is not found
	}

	return value
}

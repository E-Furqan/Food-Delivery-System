package environmentVariable

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment struct {
	HOST              string
	USER              string
	PASSWORD          string
	DB_NAME           string
	PORT              int
	JWT_SECRET        string
	RefreshTokenKey   string
	BASE_URL          string
	PROCESS_ORDER_URL string
	ORDER_PORT        string
}

func ReadEnv() Environment {
	var envVar Environment

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	envVar.HOST = os.Getenv("HOST")
	envVar.USER = os.Getenv("USER1")
	envVar.PASSWORD = os.Getenv("PASSWORD")
	envVar.DB_NAME = os.Getenv("DB_NAME")
	envVar.BASE_URL = os.Getenv("BASE_URL")
	envVar.PROCESS_ORDER_URL = os.Getenv("PROCESS_ORDER_URL")
	envVar.ORDER_PORT = os.Getenv("ORDER_PORT")

	portStr := os.Getenv("PORT")
	envVar.PORT, err = strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error converting PORT to integer: %v", err)
	}
	envVar.JWT_SECRET = os.Getenv("JWT_SECRET")
	envVar.RefreshTokenKey = os.Getenv("REFRESH_TOKEN_SECRET")
	return envVar
}

// func Get_env(key string) string {

// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}

// 	value := os.Getenv(key)
// 	if value == "" {
// 		log.Fatalf("environment variable %s is not set", key)
// 		return "" // Return an error if the variable is not found
// 	}

// 	return value
// }
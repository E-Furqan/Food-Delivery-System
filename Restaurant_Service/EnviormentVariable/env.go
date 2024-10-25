package environmentVariable

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment struct {
	HOST               string
	USER               string
	PASSWORD           string
	DB_NAME            string
	PORT               int
	JWT_SECRET         string
	RefreshTokenKey    string
	BASE_URL           string
	PROCESS_ORDER_URL  string
	GENERATE_TOKEN_URL string
	REFRESH_TOKEN_URL  string
	ORDER_PORT         string
	AUTH_PORT          string
}

var envVar Environment

func ReadEnv() Environment {

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
	envVar.AUTH_PORT = os.Getenv("AUTH_PORT")
	envVar.GENERATE_TOKEN_URL = os.Getenv("GENERATE_TOKEN_URL")
	envVar.REFRESH_TOKEN_URL = os.Getenv("REFRESH_TOKEN_URL")

	portStr := os.Getenv("PORT")
	envVar.PORT, err = strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error converting PORT to integer: %v", err)
	}
	return envVar
}

// func GetEnv(key string) string {
// 	switch key {
// 	case "HOST":
// 		return envVar.HOST
// 	case "USER1":
// 		return envVar.USER
// 	case "PASSWORD":
// 		return envVar.PASSWORD
// 	case "DB_NAME":
// 		return envVar.DB_NAME
// 	case "PORT":
// 		return strconv.Itoa(envVar.PORT)
// 	default:
// 		log.Fatalf("environment variable %s is not set", key)
// 		return ""
// 	}
// }

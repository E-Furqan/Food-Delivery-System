package environmentVariable

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment struct {
	HOST                          string
	USER                          string
	PASSWORD                      string
	DB_NAME                       string
	PORT                          int
	JWT_SECRET                    string
	RefreshTokenKey               string
	BASE_URL                      string
	PROCESS_ORDER_URL             string
	GENERATE_TOKEN_URL            string
	REFRESH_TOKEN_URL             string
	ORDER_PORT                    string
	AUTH_PORT                     string
	User_ORDERS_URL               string
	VIEW_ORDER_DETAIL_URL         string
	VIEW_ORDER_WITHOUT_DRIVER_URL string
	DRIVER_ORDERS_URL             string
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
	envVar.AUTH_PORT = os.Getenv("AUTH_PORT")
	envVar.GENERATE_TOKEN_URL = os.Getenv("GENERATE_TOKEN_URL")
	envVar.REFRESH_TOKEN_URL = os.Getenv("REFRESH_TOKEN_URL")
	envVar.User_ORDERS_URL = os.Getenv("USER_ORDERS_URL")
	envVar.DRIVER_ORDERS_URL = os.Getenv("DRIVER_ORDERS_URL")
	envVar.VIEW_ORDER_DETAIL_URL = os.Getenv("VIEW_ORDER_DETAIL_URL")
	envVar.VIEW_ORDER_WITHOUT_DRIVER_URL = os.Getenv("VIEW_ORDER_WITHOUT_DRIVER_URL")

	portStr := os.Getenv("PORT")
	envVar.PORT, err = strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error converting PORT to integer: %v", err)
	}
	envVar.JWT_SECRET = os.Getenv("JWT_SECRET")
	envVar.RefreshTokenKey = os.Getenv("REFRESH_TOKEN_SECRET")
	return envVar
}

package environmentVariable

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment struct {
	HOST                  string
	USER                  string
	PASSWORD              string
	DB_NAME               string
	PORT                  int
	JWT_SECRET            string
	RefreshTokenKey       string
	BASE_URL              string
	PROCESS_ORDER_URL     string
	GENERATE_TOKEN_URL    string
	REFRESH_TOKEN_URL     string
	ORDER_PORT            string
	AUTH_PORT             string
	RESTAURANT_ORDERS_URL string
	VIEW_ORDER_DETAIL_URL string
}

var envVar Environment

func ReadEnv() Environment {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	envVar.HOST = GetEnv("HOST", "0.0.0.0")
	envVar.USER = GetEnv("USER1", "furqan")
	envVar.PASSWORD = GetEnv("PASSWORD", "furqan")
	envVar.DB_NAME = GetEnv("DB_NAME", "Restaurant")
	envVar.BASE_URL = GetEnv("BASE_URL", "http://localhost")
	envVar.PROCESS_ORDER_URL = GetEnv("PROCESS_ORDER_URL", "/order/update/status")
	envVar.ORDER_PORT = GetEnv("ORDER_PORT", ":8081")
	envVar.AUTH_PORT = GetEnv("AUTH_PORT", ":8084")
	envVar.GENERATE_TOKEN_URL = GetEnv("GENERATE_TOKEN_URL", "/auth/login")
	envVar.REFRESH_TOKEN_URL = GetEnv("REFRESH_TOKEN_URL", "Ali")
	envVar.RESTAURANT_ORDERS_URL = GetEnv("RESTAURANT_ORDERS_URL", "/order/view/restaurant/orders")
	envVar.VIEW_ORDER_DETAIL_URL = GetEnv("VIEW_ORDER_DETAIL_URL", "/order/view/order")
	envVar.JWT_SECRET = GetEnv("JWT_SECRET", "Furqan")
	portStr := GetEnv("PORT", "5432")
	envVar.PORT, err = strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error converting PORT to integer: %v", err)
		envVar.PORT = 5432
	}
	return envVar
}
func GetEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

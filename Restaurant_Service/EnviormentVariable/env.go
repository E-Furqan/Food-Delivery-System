package environmentVariable

import (
	"log"
	"strconv"

	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/joho/godotenv"
)

type Environment struct {
	DATABASE_HOST           string
	DATABASE_USER           string
	DATABASE_PASSWORD       string
	DATABASE_NAME           string
	DATABASE_PORT           int
	JWT_SECRET              string
	RefreshTokenKey         string
	BASE_URL                string
	UPDATE_ORDER_STATUS_URL string
	GENERATE_TOKEN_URL      string
	REFRESH_TOKEN_URL       string
	ORDER_PORT              string
	AUTH_PORT               string
	RESTAURANT_ORDERS_URL   string
	VIEW_ORDER_DETAIL_URL   string
}

var envVar Environment

func ReadEnv() Environment {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	envVar.DATABASE_HOST = utils.GetEnv("DATABASE_HOST", "db")
	envVar.DATABASE_USER = utils.GetEnv("DATABASE_USER", "furqan")
	envVar.DATABASE_PASSWORD = utils.GetEnv("DATABASE_PASSWORD", "furqan")
	envVar.DATABASE_NAME = utils.GetEnv("DATABASE_NAME", "Restaurant")
	envVar.BASE_URL = utils.GetEnv("BASE_URL", "http://localhost")
	envVar.UPDATE_ORDER_STATUS_URL = utils.GetEnv("UPDATE_ORDER_STATUS_URL", "/order/update/status")
	envVar.ORDER_PORT = utils.GetEnv("ORDER_PORT", ":8081")
	envVar.AUTH_PORT = utils.GetEnv("AUTH_PORT", ":8084")
	envVar.GENERATE_TOKEN_URL = utils.GetEnv("GENERATE_TOKEN_URL", "/auth/generate/token")
	envVar.REFRESH_TOKEN_URL = utils.GetEnv("REFRESH_TOKEN_URL", "/auth/refresh/token")
	envVar.RESTAURANT_ORDERS_URL = utils.GetEnv("RESTAURANT_ORDERS_URL", "/order/view/restaurant/orders")
	envVar.VIEW_ORDER_DETAIL_URL = utils.GetEnv("VIEW_ORDER_DETAIL_URL", "/order/view/order")
	envVar.JWT_SECRET = utils.GetEnv("JWT_SECRET", "Furqan")
	portStr := utils.GetEnv("PORT", "5432")
	envVar.DATABASE_PORT, err = strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error converting PORT to integer: %v", err)
		envVar.DATABASE_PORT = 5432
	}
	return envVar
}

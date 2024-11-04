package environmentVariable

import (
	"log"
	"strconv"

	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/joho/godotenv"
)

type Environment struct {
	HOST                    string
	USER                    string
	PASSWORD                string
	DB_NAME                 string
	PORT                    int
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

	envVar.HOST = utils.GetEnv("HOST", "0.0.0.0")
	envVar.USER = utils.GetEnv("USER1", "furqan")
	envVar.PASSWORD = utils.GetEnv("PASSWORD", "furqan")
	envVar.DB_NAME = utils.GetEnv("DB_NAME", "Restaurant")
	envVar.BASE_URL = utils.GetEnv("BASE_URL", "http://localhost")
	envVar.UPDATE_ORDER_STATUS_URL = utils.GetEnv("UPDATE_ORDER_STATUS_URL", "/order/update/status")
	envVar.ORDER_PORT = utils.GetEnv("ORDER_PORT", ":8081")
	envVar.AUTH_PORT = utils.GetEnv("AUTH_PORT", ":8084")
	envVar.GENERATE_TOKEN_URL = utils.GetEnv("GENERATE_TOKEN_URL", "/auth/login")
	envVar.REFRESH_TOKEN_URL = utils.GetEnv("REFRESH_TOKEN_URL", "Ali")
	envVar.RESTAURANT_ORDERS_URL = utils.GetEnv("RESTAURANT_ORDERS_URL", "/order/view/restaurant/orders")
	envVar.VIEW_ORDER_DETAIL_URL = utils.GetEnv("VIEW_ORDER_DETAIL_URL", "/order/view/order")
	envVar.JWT_SECRET = utils.GetEnv("JWT_SECRET", "Furqan")
	portStr := utils.GetEnv("PORT", "5432")
	envVar.PORT, err = strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error converting PORT to integer: %v", err)
		envVar.PORT = 5432
	}
	return envVar
}

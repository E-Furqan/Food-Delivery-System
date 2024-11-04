package environmentVariable

import (
	"log"
	"strconv"

	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
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
	UPDATE_ORDER_STATUS_URL       string
	GENERATE_TOKEN_URL            string
	REFRESH_TOKEN_URL             string
	ORDER_PORT                    string
	AUTH_PORT                     string
	USER_ORDERS_URL               string
	VIEW_ORDER_DETAIL_URL         string
	VIEW_ORDER_WITHOUT_DRIVER_URL string
	DRIVER_ORDERS_URL             string
	ASSIGN_DRIVER_URL             string
}

func ReadEnv() Environment {
	var envVar Environment

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	envVar.HOST = utils.GetEnv("HOST", "0.0.0.0")
	envVar.USER = utils.GetEnv("USER1", "furqan")
	envVar.PASSWORD = utils.GetEnv("PASSWORD", "furqan")
	envVar.DB_NAME = utils.GetEnv("DB_NAME", "User")

	envVar.BASE_URL = utils.GetEnv("BASE_URL", "http://localhost")
	envVar.UPDATE_ORDER_STATUS_URL = utils.GetEnv("UPDATE_ORDER_STATUS_URL", "/order/update/status")
	envVar.GENERATE_TOKEN_URL = utils.GetEnv("GENERATE_TOKEN_URL", "/auth/generate/token")
	envVar.REFRESH_TOKEN_URL = utils.GetEnv("REFRESH_TOKEN_URL", "/auth/refresh/token")
	envVar.USER_ORDERS_URL = utils.GetEnv("USER_ORDERS_URL", "/order/view/user/orders")
	envVar.DRIVER_ORDERS_URL = utils.GetEnv("DRIVER_ORDERS_URL", "/order/view/driver/orders")
	envVar.VIEW_ORDER_DETAIL_URL = utils.GetEnv("VIEW_ORDER_DETAIL_URL", "/order/view/order")
	envVar.VIEW_ORDER_WITHOUT_DRIVER_URL = utils.GetEnv("VIEW_ORDER_WITHOUT_DRIVER_URL", "/order/view/without/driver/orders")
	envVar.ASSIGN_DRIVER_URL = utils.GetEnv("ASSIGN_DRIVER_URL", "/order/assign/diver")

	envVar.ORDER_PORT = utils.GetEnv("ORDER_PORT", ":8081")
	envVar.AUTH_PORT = utils.GetEnv("AUTH_PORT", ":8084")
	portStr := utils.GetEnv("PORT", "5430")
	envVar.PORT, err = strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error converting PORT to integer: %v", err)
		envVar.PORT = 5430
	}
	envVar.JWT_SECRET = utils.GetEnv("JWT_SECRET", "Furqan")
	envVar.RefreshTokenKey = utils.GetEnv("REFRESH_TOKEN_SECRET", "Ali")
	return envVar
}

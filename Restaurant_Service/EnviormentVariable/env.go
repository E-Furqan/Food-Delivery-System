package environmentVariable

import (
	"log"
	"strconv"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/joho/godotenv"
)

func ReadDatabaseEnv() model.DatabaseEnv {
	var envVar model.DatabaseEnv
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	envVar.DATABASE_HOST = utils.GetEnv("DATABASE_HOST", "db")
	envVar.DATABASE_USER = utils.GetEnv("DATABASE_USER", "furqan")
	envVar.DATABASE_PASSWORD = utils.GetEnv("DATABASE_PASSWORD", "furqan")
	envVar.DATABASE_NAME = utils.GetEnv("DATABASE_NAME", "Restaurant")
	portStr := utils.GetEnv("PORT", "5432")
	envVar.DATABASE_PORT, err = strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Error converting PORT to integer: %v", err)
		envVar.DATABASE_PORT = 5432
	}
	return envVar
}

func ReadOrderClientEnv() model.OrderClientEnv {
	var envVar model.OrderClientEnv
	envVar.BASE_URL = utils.GetEnv("BASE_URL", "http://localhost")
	envVar.UPDATE_ORDER_STATUS_URL = utils.GetEnv("UPDATE_ORDER_STATUS_URL", "/order/update/status")
	envVar.ORDER_PORT = utils.GetEnv("ORDER_PORT", ":8081")
	envVar.RESTAURANT_ORDERS_URL = utils.GetEnv("RESTAURANT_ORDERS_URL", "/order/view/orders")
	envVar.VIEW_ORDER_DETAIL_URL = utils.GetEnv("VIEW_ORDER_DETAIL_URL", "/order/view/order")
	return envVar
}

func ReadAuthClientEnv() model.AuthClientEnv {
	var envVar model.AuthClientEnv
	envVar.AUTH_PORT = utils.GetEnv("AUTH_PORT", ":8084")
	envVar.GENERATE_TOKEN_URL = utils.GetEnv("GENERATE_TOKEN_URL", "/auth/generate/token")
	envVar.REFRESH_TOKEN_URL = utils.GetEnv("REFRESH_TOKEN_URL", "/auth/refresh/token")
	return envVar
}

func ReadMiddlewareEnv() model.MiddlewareEnv {
	var envVar model.MiddlewareEnv
	envVar.JWT_SECRET = utils.GetEnv("JWT_SECRET", "Furqan")
	envVar.REFRESH_TOKEN_SECRET = utils.GetEnv("REFRESH_TOKEN_SECRET", "Ali")
	return envVar
}

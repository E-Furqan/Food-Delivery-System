package environmentVariable

import (
	"log"
	"strconv"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/joho/godotenv"
)

func ReadDatabaseEnv() model.DatabaseEnv {
	var DatabaseEnv model.DatabaseEnv

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	DatabaseEnv.DATABASE_HOST = utils.GetEnv("DATABASE_HOST", "db")
	DatabaseEnv.DATABASE_USER = utils.GetEnv("DATABASE_USER", "furqan")
	DatabaseEnv.DATABASE_PASSWORD = utils.GetEnv("DATABASE_PASSWORD", "furqan")
	DatabaseEnv.DATABASE_NAME = utils.GetEnv("DATABASE_NAME", "User")

	portStr := utils.GetEnv("DATABASE_PORT", "5432")
	DatabaseEnv.DATABASE_PORT, err = strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Error converting PORT to integer: %v", err)
		DatabaseEnv.DATABASE_PORT = 5432
	}

	return DatabaseEnv
}

func ReadOrderClientEnv() model.OrderClientEnv {
	var OrderClientEnv model.OrderClientEnv

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	OrderClientEnv.BASE_URL = utils.GetEnv("BASE_URL", "http://localhost")
	OrderClientEnv.UPDATE_ORDER_STATUS_URL = utils.GetEnv("UPDATE_ORDER_STATUS_URL", "/order/update/status")
	OrderClientEnv.VIEW_ORDERS_URL = utils.GetEnv("USER_ORDERS_URL", "/order/view/orders")
	OrderClientEnv.VIEW_ORDER_WITHOUT_DRIVER_URL = utils.GetEnv("VIEW_ORDER_WITHOUT_DRIVER_URL", "/order/view/without/driver/orders")
	OrderClientEnv.ASSIGN_DRIVER_URL = utils.GetEnv("ASSIGN_DRIVER_URL", "/order/assign/diver")
	OrderClientEnv.ORDER_PORT = utils.GetEnv("ORDER_PORT", ":8081")

	return OrderClientEnv
}

func ReadAuthClientEnv() model.AuthClientEnv {
	var AuthClientEnv model.AuthClientEnv

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	AuthClientEnv.BASE_URL = utils.GetEnv("BASE_URL", "http://localhost")

	AuthClientEnv.GENERATE_TOKEN_URL = utils.GetEnv("GENERATE_TOKEN_URL", "/auth/generate/token")
	AuthClientEnv.REFRESH_TOKEN_URL = utils.GetEnv("REFRESH_TOKEN_URL", "/auth/refresh/token")

	AuthClientEnv.AUTH_PORT = utils.GetEnv("AUTH_PORT", ":8084")
	return AuthClientEnv
}

func ReadMiddlewareEnv() model.MiddlewareEnv {
	var MiddlewareEnv model.MiddlewareEnv

	MiddlewareEnv.JWT_SECRET = utils.GetEnv("JWT_SECRET", "Furqan")
	MiddlewareEnv.RefreshTokenKey = utils.GetEnv("REFRESH_TOKEN_SECRET", "Ali")

	return MiddlewareEnv
}

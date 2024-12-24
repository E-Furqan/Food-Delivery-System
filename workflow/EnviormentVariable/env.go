package environmentVariable

import (
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/joho/godotenv"
)

func ReadOrderClientEnv() model.OrderClientEnv {
	var OrderClientEnv model.OrderClientEnv

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	OrderClientEnv.BASE_URL = utils.GetEnv("BASE_URL", "http://localhost")
	OrderClientEnv.UPDATE_ORDER_STATUS_URL = utils.GetEnv("UPDATE_ORDER_STATUS_URL", "/order/update/status")
	OrderClientEnv.Fetch_OrderStatus_URL = utils.GetEnv("UPDATE_ORDER_STATUS_URL", "/order/fetch/order/status")
	OrderClientEnv.Create_Order_URL = utils.GetEnv("UPDATE_ORDER_STATUS_URL", "/order/create/order")
	OrderClientEnv.ORDER_PORT = utils.GetEnv("ORDER_PORT", "8081")

	return OrderClientEnv
}

func ReadRestaurantClientEnv() model.RestaurantClientEnv {
	var envVar model.RestaurantClientEnv

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	envVar.BASE_URL = utils.GetEnv("BASE_URL", "http://localhost")
	envVar.Get_Items_URL = utils.GetEnv("GET_ITEMS_URL", "/restaurant/fetch/item/prices")
	envVar.RESTAURANT_PORT = utils.GetEnv("RESTAURANT_PORT", ":8082")

	return envVar
}

func ReadEmailClientEnv() model.EmailEnv {
	var envVar model.EmailEnv

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	envVar.EmailAddressFrom = utils.GetEnv("EMAIL_FROM", "furqan.ali@emumba.com")
	envVar.EmailPassKey = utils.GetEnv("PASS_KEY", "sqrf gefw qccw pqyr")
	return envVar
}

func ReadUserClientEnv() model.UserClientEnv {
	var envVar model.UserClientEnv

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	envVar.BASE_URL = utils.GetEnv("BASE_URL", "http://localhost")
	envVar.Fetch_email_URL = utils.GetEnv("GET_ITEMS_URL", "/user/fetch/email")
	envVar.USER_PORT = utils.GetEnv("RESTAURANT_PORT", "8083")

	return envVar
}

func ReadPipelineEnv() model.DatapipelineClientEnv {
	var envVar model.DatapipelineClientEnv

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	envVar.BASE_URL = utils.GetEnv("BASE_URL", "http://localhost")
	envVar.FETCH_SOURCE_CONFIGURATION_URL = utils.GetEnv("FETCH_SOURCE_CONFIGURATION_URL", "/pipeline/fetch/source/configuration")
	envVar.FETCH_DESTINATION_CONFIGURATION_URL = utils.GetEnv("FETCH_DESTINATION_CONFIGURATION_URL", "/pipeline/fetch/destination/configuration")
	envVar.ADD_LOGS_URL = utils.GetEnv("ADD_LOGS_URL", "/pipeline/add/logs")
	envVar.DATAPIPELINE_PORT = utils.GetEnv("DATAPIPELINE_PORT", "8085")

	return envVar
}

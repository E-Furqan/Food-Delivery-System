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

	DatabaseEnv.DATABASE_HOST = utils.GetEnv("DATABASE_HOST", "localhost")
	DatabaseEnv.DATABASE_USER = utils.GetEnv("DATABASE_USER", "furqan")
	DatabaseEnv.DATABASE_PASSWORD = utils.GetEnv("DATABASE_PASSWORD", "furqan")
	DatabaseEnv.DATABASE_NAME = utils.GetEnv("DATABASE_NAME", "DataPipeline")

	portStr := utils.GetEnv("DATABASE_PORT", "5435")
	DatabaseEnv.DATABASE_PORT, err = strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Error converting PORT to integer: %v", err)
		DatabaseEnv.DATABASE_PORT = 5432
	}

	return DatabaseEnv
}

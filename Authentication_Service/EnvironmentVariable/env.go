package environmentVariable

import (
	"log"

	payload "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Payload"
	utils "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Utils"
	"github.com/joho/godotenv"
)

var envVar payload.Environment

func ReadEnv() payload.Environment {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	envVar.JWT_SECRET = utils.GetEnv("JWT_SECRET", "Furqan")
	envVar.RefreshTokenKey = utils.GetEnv("RefreshTokenKey", "Ali")
	return envVar
}

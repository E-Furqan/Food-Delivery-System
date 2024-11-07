package environmentVariable

import (
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Model"

	utils "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Utils"
	"github.com/joho/godotenv"
)

var envVar model.AuthSecrets

func ReadEnv() model.AuthSecrets {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	envVar.JWT_SECRET = utils.GetEnv("JWT_SECRET", "Furqan")
	envVar.RefreshTokenKey = utils.GetEnv("RefreshTokenKey", "Ali")
	return envVar
}

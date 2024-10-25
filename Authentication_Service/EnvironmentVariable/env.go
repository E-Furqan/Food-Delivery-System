package environmentVariable

import (
	"log"
	"os"

	payload "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Payload"
	"github.com/joho/godotenv"
)

var envVar payload.Environment

func ReadEnv() payload.Environment {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	envVar.JWT_SECRET = os.Getenv("JWT_SECRET")
	envVar.RefreshTokenKey = os.Getenv("RefreshTokenKey")
	return envVar
}

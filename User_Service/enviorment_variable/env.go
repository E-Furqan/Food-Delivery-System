package environmentvariable

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	HOST       string
	USER       string
	PASSWORD   string
	DB_NAME    string
	PORT       int
	JWT_SECRET string
	ADMIN      string
	ADMIN_PASS string
}

// // ReadEnv reads environment variables from a .env file and returns an Environment struct
// func ReadEnv() Environment {
// 	var envVar Environment

// 	// Load the .env file
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}

// 	// Read environment variables
// 	envVar.HOST = os.Getenv("HOST")
// 	envVar.USER = os.Getenv("USER1")
// 	envVar.PASSWORD = os.Getenv("PASSWORD")
// 	envVar.DB_NAME = os.Getenv("DB_NAME")
// 	portStr := os.Getenv("PORT")
// 	envVar.PORT, err = strconv.Atoi(portStr)
// 	if err != nil {
// 		log.Fatalf("Error converting PORT to integer: %v", err)
// 	}
// 	envVar.JWT_SECRET = os.Getenv("JWT_SECRET")
// 	envVar.ADMIN = os.Getenv("ADMIN")
// 	envVar.ADMIN_PASS = os.Getenv("ADMIN_PASS")

// 	return envVar
// }

func Get_env(key string) string {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return os.Getenv(key)
}

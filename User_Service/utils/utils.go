package utils

import (
	"log"
	"time"

	environmentvariable "github.com/E-Furqan/Food-Delivery-System/enviorment_variable"
	"github.com/dgrijalva/jwt-go"
)

var env = environmentvariable.ReadEnv()
var jwtKey = []byte(env.JWT_SECRET)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generating token: %v", err) // Log error
		return "", err
	}

	return tokenString, nil
}

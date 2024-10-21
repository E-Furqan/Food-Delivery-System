package utils

import (
	"log"
	"net/http"
	"time"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey []byte
var refreshTokenKey []byte

// add active token and when user switch role generate a new jwt token
type Claims struct {
	Username string `json:"username"`
	RoleId   []uint `json:"RoleId"`
	jwt.StandardClaims
}

func SetEnvValue(envVar environmentVariable.Environment) {
	jwtKey = []byte(envVar.JWT_SECRET)
	refreshTokenKey = []byte(envVar.RefreshTokenKey)
}

func GenerateTokens(username string, roleId []uint) (string, string, error) {

	accessExpirationTime := time.Now().Add(15 * time.Minute)
	accessClaims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpirationTime.Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return "", "", err
	}

	refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 days
	refreshClaims := &Claims{
		Username: username,
		RoleId:   roleId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(refreshTokenKey)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func RefreshToken(refreshToken string, c *gin.Context) (string, error) {

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshTokenKey), nil
	})

	if err != nil {

		log.Printf("Error parsing refresh token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return "", err
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return "", nil
	}

	accessToken, _, err := GenerateTokens(claims.Username, claims.RoleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new access token"})
		return "", err
	}

	return accessToken, nil
}

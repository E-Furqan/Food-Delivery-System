package utils

import (
	"log"
	"net/http"
	"os"
	"time"

	model "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtKey []byte
var refreshTokenKey []byte

func SetEnvValue(envVar model.AuthSecrets) {
	jwtKey = []byte(envVar.JWT_SECRET)
	refreshTokenKey = []byte(envVar.RefreshTokenKey)
}

func GenerateTokens(accessClaims model.Claims, refreshClaims model.Claims) (string, string, error) {

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return "", "", err
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

	claims := &model.Input{}
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
	var accessClaims model.Claims
	var refreshClaims model.Claims

	log.Print("claims.id")
	log.Print(claims.ClaimId)
	log.Print("claims.role")
	log.Print(claims.ActiveRole)

	var input model.Input
	input.ClaimId = claims.ClaimId
	input.ActiveRole = claims.ActiveRole
	accessClaims, refreshClaims = CreateClaim(input)

	accessToken, _, err := GenerateTokens(accessClaims, refreshClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new access token"})
		return "", err
	}

	return accessToken, nil
}

func CreateClaim(input model.Input) (model.Claims, model.Claims) {
	var accessClaims model.Claims
	var refreshClaims model.Claims

	accessClaims = &model.GeneralClaim{
		ClaimId:    input.ClaimId,
		ActiveRole: input.ActiveRole,
	}
	accessClaims.SetExpirationTime(time.Now().Add(15 * time.Minute).Unix())

	refreshClaims = &model.GeneralClaim{
		ClaimId:    input.ClaimId,
		ActiveRole: input.ActiveRole,
	}

	refreshClaims.SetExpirationTime(time.Now().Add(7 * 24 * time.Hour).Unix())
	return accessClaims, refreshClaims
}

func GetEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

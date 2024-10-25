package utils

import (
	"log"
	"net/http"
	"time"

	payload "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Payload"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtKey []byte
var refreshTokenKey []byte

func SetEnvValue(envVar payload.Environment) {
	jwtKey = []byte(envVar.JWT_SECRET)
	refreshTokenKey = []byte(envVar.RefreshTokenKey)
}

func GenerateTokens(accessClaims payload.Claims, refreshClaims payload.Claims) (string, string, error) {

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

type Claims struct {
	// RestaurantID uint   `json:"restaurant_id"`
	// OrderId      uint   `json:"order_id"`
	ClaimId     uint   `json:"claim_id"`
	Username    string `json:"username"`
	ActiveRole  string `json:"activeRole"`
	ServiceType string `json:"service_type"`
	jwt.StandardClaims
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
	var accessClaims payload.Claims
	var refreshClaims payload.Claims

	log.Print("claims.RestaurantID")
	log.Print(claims.ClaimId)
	log.Print("claims.ServiceType")
	log.Print(claims.ServiceType)

	var input payload.Input

	if claims.ServiceType == "User" {
		input.Username = claims.Username
		input.ActiveRole = claims.ActiveRole
		accessClaims, refreshClaims = CreateUserClaim(input)
	} else {
		input.ClaimId = claims.ClaimId
		accessClaims, refreshClaims = CreateClaim(input)
	}

	accessToken, _, err := GenerateTokens(accessClaims, refreshClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new access token"})
		return "", err
	}

	return accessToken, nil
}

func CreateUserClaim(input payload.Input) (payload.Claims, payload.Claims) {
	var accessClaims payload.Claims
	var refreshClaims payload.Claims

	accessClaims = &payload.UserClaims{
		Username:    input.Username,
		ActiveRole:  input.ActiveRole,
		ServiceType: input.ServiceType,
	}
	accessClaims.SetExpirationTime(time.Now().Add(30 * time.Minute).Unix())

	refreshClaims = &payload.UserClaims{
		Username:    input.Username,
		ActiveRole:  input.ActiveRole,
		ServiceType: "User",
	}
	refreshClaims.SetExpirationTime(time.Now().Add(7 * 24 * time.Hour).Unix())

	return accessClaims, refreshClaims
}

func CreateClaim(input payload.Input) (payload.Claims, payload.Claims) {
	var accessClaims payload.Claims
	var refreshClaims payload.Claims

	accessClaims = &payload.RestaurantClaims{
		ClaimId:     input.ClaimId,
		ServiceType: input.ServiceType,
	}
	accessClaims.SetExpirationTime(time.Now().Add(15 * time.Minute).Unix())

	refreshClaims = &payload.RestaurantClaims{
		ClaimId:     input.ClaimId,
		ServiceType: input.ServiceType,
	}

	refreshClaims.SetExpirationTime(time.Now().Add(7 * 24 * time.Hour).Unix())
	return accessClaims, refreshClaims
}

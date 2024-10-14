package utils

import (
	"log"
	"net/http"
	"time"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/enviorment_variable"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte(environmentVariable.Get_env("JWT_SECRET"))
var refreshTokenKey = []byte(environmentVariable.Get_env("REFRESH_TOKEN_SECRET"))

type Claims struct {
	Username string `json:"username"`
	RoleId   []uint `json:"RoleId"`
	jwt.StandardClaims
}

// auth
func GenerateTokens(username string, roleId []uint) (string, string, error) {
	// Access Token expiration time
	accessExpirationTime := time.Now().Add(15 * time.Minute) // 15 minutes
	accessClaims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpirationTime.Unix(),
		},
	}

	// Generate Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return "", "", err
	}

	// Refresh Token expiration time
	refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 days
	refreshClaims := &Claims{
		Username: username,
		RoleId:   roleId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
		},
	}

	// Generate Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(refreshTokenKey)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// RefreshToken generates a new access token using a valid refresh token
func RefreshToken(refreshToken string, c *gin.Context) (string, error) {
	// Parse the refresh token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshTokenKey), nil // Ensure you are returning the key as a byte slice
	})

	if err != nil {
		// Log the error for debugging
		log.Printf("Error parsing refresh token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return "", err
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return "", nil // Returning nil error for invalid token (client-side handling)
	}

	// Generate new access token
	accessToken, _, err := GenerateTokens(claims.Username, claims.RoleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new access token"})
		return "", err
	}

	return accessToken, nil // Return the new access token and no error
}

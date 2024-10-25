package utils

import (
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	payload.RestaurantClaim
	jwt.StandardClaims
}

// func GenerateTokens(RestaurantId uint) (string, string, error) {

// 	accessExpirationTime := time.Now().Add(15 * time.Minute)
// 	accessClaims := &Claims{
// 		RestaurantID: RestaurantId,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: accessExpirationTime.Unix(),
// 		},
// 	}

// 	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
// 	accessTokenString, err := accessToken.SignedString(jwtKey)
// 	if err != nil {
// 		log.Printf("Error generating access token: %v", err)
// 		return "", "", err
// 	}

// 	refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 days
// 	refreshClaims := &Claims{
// 		RestaurantID: RestaurantId,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: refreshExpirationTime.Unix(),
// 		},
// 	}

// 	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
// 	refreshTokenString, err := refreshToken.SignedString(refreshTokenKey)
// 	if err != nil {
// 		log.Printf("Error generating refresh token: %v", err)
// 		return "", "", err
// 	}

// 	return accessTokenString, refreshTokenString, nil
// }

// func RefreshToken(refreshToken string, c *gin.Context) (string, error) {

// 	claims := &Claims{}
// 	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(refreshTokenKey), nil
// 	})
// 	log.Print(claims.RestaurantID)
// 	log.Print(claims.ServiceType)
// 	if err != nil {

// 		log.Printf("Error parsing refresh token: %v", err)
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
// 		return "", err
// 	}

// 	if !token.Valid {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
// 		return "", nil
// 	}

// 	accessToken, _, err := GenerateTokens(claims.RestaurantID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new access token"})
// 		return "", err
// 	}

// 	return accessToken, nil
// }

func GenerateResponse(httpStatusCode int, c *gin.Context, title1 string, message1 string, title2 string, input interface{}) {
	response := gin.H{
		title1: message1,
	}

	if title2 != "" && input != nil {
		response[title2] = input
	}

	c.JSON(httpStatusCode, response)
}

package utils

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
<<<<<<< HEAD:Authentication_Service/Utils/utils.go
	"os"
	"time"

	model "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Model"
=======
	"net/url"
	"os"
	"strings"

>>>>>>> b0a1439a54cf96a16dcaeb351bfed09905ab01c9:Restaurant_Service/Utils/utils.go
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

<<<<<<< HEAD:Authentication_Service/Utils/utils.go
var jwtKey []byte
var refreshTokenKey []byte

func SetEnvValue(envVar model.AuthSecrets) {
	jwtKey = []byte(envVar.JWT_SECRET)
	refreshTokenKey = []byte(envVar.RefreshTokenKey)
}

func GenerateTokens(accessClaims model.Claims, refreshClaims model.Claims) (string, string, error) {
=======
func Verification(c *gin.Context) (any, error) {
	RestaurantID, exists := c.Get("RestaurantID")
	if !exists {
		return 0, fmt.Errorf("restaurant id does not exist")
	}
>>>>>>> b0a1439a54cf96a16dcaeb351bfed09905ab01c9:Restaurant_Service/Utils/utils.go

	RestaurantID, ok := RestaurantID.(uint)
	if !ok {
		return 0, fmt.Errorf("invalid restaurant id")
	}
<<<<<<< HEAD:Authentication_Service/Utils/utils.go

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
=======
	return RestaurantID, nil
>>>>>>> b0a1439a54cf96a16dcaeb351bfed09905ab01c9:Restaurant_Service/Utils/utils.go
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

func GetEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func GetAuthToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization token not provided")
	}
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}
	token := tokenParts[1]
	return token, nil
}

func CreateAuthorizedRequest(url string, jsonData []byte, c *gin.Context, MethodType string) (*http.Request, error) {

	req, err := http.NewRequest(MethodType, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	token, err := GetAuthToken(c)
	if err != nil {
		GenerateResponse(http.StatusUnauthorized, c, "error", err.Error(), "", nil)
		return nil, fmt.Errorf("error retrieving auth token of user: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return req, nil
}

func CreateRequest(url string, jsonData []byte, MethodType string) (*http.Request, error) {
	req, err := http.NewRequest(MethodType, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func CreateHTTPClient() *http.Client {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return client
}

func CreateUrl(BaseUrl string, Port string, APIUrl string) (string, error) {

	if !strings.HasPrefix(BaseUrl, "http://") && !strings.HasPrefix(BaseUrl, "https://") {
		return "", errors.New("BaseUrl must start with http:// or https://")
	}

	if _, err := url.ParseRequestURI(Port); err != nil {
		return "", fmt.Errorf("invalid Port: %v", err)
	}

	baseURL, err := url.Parse(BaseUrl)
	if err != nil {
		return "", fmt.Errorf("invalid BaseUrl: %v", err)
	}
	baseURL.Host = fmt.Sprintf("%s:%s", baseURL.Hostname(), Port)

	escapedAPIUrl := url.PathEscape(APIUrl)

	finalURL := fmt.Sprintf("%s%s", baseURL.String(), escapedAPIUrl)
	return finalURL, nil
}

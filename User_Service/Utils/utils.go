package utils

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/gin-gonic/gin"
)

func GenerateResponse(httpStatusCode int, c *gin.Context, title1 string, message1 string, title2 string, input interface{}) {

	errorMessage := strings.TrimPrefix(message1, "ERROR: ")
	response := gin.H{
		title1: errorMessage,
	}

	if title2 != "" && input != nil {
		response[title2] = input
	}

	c.JSON(httpStatusCode, response)
}

func VerifyUserId(c *gin.Context) (any, error) {
	userIdValue, exists := c.Get("userId")
	if !exists {
		return 0, fmt.Errorf("userId does not exist")
	}
	userId, ok := userIdValue.(uint)
	if !ok {
		return 0, fmt.Errorf("invalid userId")
	}
	return userId, nil
}

func VerifyActiveAdminRole(c *gin.Context) (any, error) {
	activeRole, err := FetchActiveRole(c)
	if err != nil {
		return activeRole, err
	}

	if activeRole != "Admin" {
		return activeRole, fmt.Errorf("insufficient permissions")
	}

	return activeRole, nil
}

func CreateUserClaim(user model.User) model.UserClaim {
	var UserClaim model.UserClaim
	UserClaim.UserId = user.UserId
	UserClaim.ActiveRole = user.ActiveRole
	return UserClaim
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

func FetchActiveRole(c *gin.Context) (any, error) {

	activeRole, exists := c.Get("activeRole")
	if !exists {
		return nil, fmt.Errorf("user role does not exist")
	}

	return activeRole, nil
}

func VerifyIfDriver(activeRole any) error {

	if activeRole != "Delivery driver" {
		return fmt.Errorf("insufficient permission")
	}

	return nil
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

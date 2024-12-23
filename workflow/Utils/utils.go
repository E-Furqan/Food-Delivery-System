package utils

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/api/drive/v3"
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

func FetchClaimsUserId(c *gin.Context) (any, error) {
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

func CreateAuthorizedRequest(url string, jsonData []byte, MethodType string, token string) (*http.Request, error) {

	req, err := http.NewRequest(MethodType, url, bytes.NewBuffer(jsonData))
	log.Print("create order url:", url, MethodType)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
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

	finalURL := fmt.Sprintf("%s:%s%s", BaseUrl, Port, APIUrl)
	return finalURL, nil
}

func EmailGenerator(orderID uint, orderStatus string) ([]byte, error) {
	var subject, body string

	switch strings.ToLower(orderStatus) {
	case OrderPlaced:
		subject = OrderPlacedSubject
		body = fmt.Sprintf("Great news! Your order with Order ID %v has been successfully placed. Please wait for the restaurant to confirm your order. Thank you for choosing us.", orderID)

	case Cancelled:
		subject = CancelledSubject
		body = fmt.Sprintf("We regret to inform you that your order with Order ID %v has been cancelled. If you have any questions, please contact support.", orderID)

	case Accepted:
		subject = AcceptedSubject
		body = fmt.Sprintf("Great news! Your order with Order ID %v has been confirmed. Thank you for choosing us.", orderID)

	case Completed:
		subject = CompletedSubject
		body = fmt.Sprintf("Congratulations! Your order with Order ID %v has been successfully completed. We would appreciate it if you could leave a review. Thank you!", orderID)

	case Delay:
		subject = DelaySubject
		body = fmt.Sprintf("We regret to inform you that your order with Order ID %v has been delayed because the restaurant has not accepted it yet.", orderID)

	default:
		log.Println("Invalid order status provided.")
		return []byte{}, fmt.Errorf("invalid order status: %s", orderStatus)
	}

	message := []byte(subject + "\n" + body)
	return message, nil
}

func ActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 10,
			MaximumInterval:    time.Second * 30,
			MaximumAttempts:    3,
			BackoffCoefficient: 2.0,
		},
	}
}

func UpdateOrderStatusTOCancel(orderID uint) model.UpdateOrder {

	var updatedOrder model.UpdateOrder
	updatedOrder.OrderId = orderID
	updatedOrder.OrderStatus = Cancelled
	updatedOrder.UserID = 0
	return updatedOrder
}

func Sleep(ctx workflow.Context) {
	workflow.Sleep(ctx, 2*time.Second)
}

func ListFilesInFolder(client *drive.Service, folderID string) ([]*drive.File, error) {
	query := fmt.Sprintf("'%s' in parents and trashed = false", folderID)
	fileList, err := client.Files.List().Q(query).Do()
	if err != nil {
		return nil, err
	}
	return fileList.Files, nil
}

func ExtractFolderID(folderUrl string) (string, error) {

	parts := strings.Split(folderUrl, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid folder URL")
	}
	return parts[len(parts)-1], nil
}

func CreateSourceObj(SourcesID int) model.Source {
	var source model.Source

	source.SourcesID = SourcesID
	return source
}

func CreateDestinationObj(SourcesID int) model.Destination {
	var Destination model.Destination

	Destination.DestinationsID = SourcesID
	return Destination
}

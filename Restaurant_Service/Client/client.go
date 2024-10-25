package ClientPackage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
)

type Client struct {
	BaseUrl               string
	ProcessOrderURL       string
	OrderPort             string
	AuthPort              string
	GenerateResponseUrl   string
	RefreshTokenUrl       string
	RESTAURANT_ORDERS_URL string
	VIEW_ORDER_DETAIL_URL string
}

func NewClient() *Client {
	return &Client{}
}

func (client *Client) SetEnvValue(envVar environmentVariable.Environment) {
	client.BaseUrl = envVar.BASE_URL
	client.ProcessOrderURL = envVar.PROCESS_ORDER_URL
	client.OrderPort = envVar.ORDER_PORT
	client.AuthPort = envVar.AUTH_PORT
	client.GenerateResponseUrl = envVar.GENERATE_TOKEN_URL
	client.RefreshTokenUrl = envVar.REFRESH_TOKEN_URL
	client.RESTAURANT_ORDERS_URL = envVar.RESTAURANT_ORDERS_URL
	client.VIEW_ORDER_DETAIL_URL = envVar.VIEW_ORDER_DETAIL_URL

}

func (client *Client) ProcessOrder(input payload.ProcessOrder) error {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("error marshaling input: %v", err)
	}

	url := fmt.Sprintf("%s%s%s", client.BaseUrl, client.OrderPort, client.ProcessOrderURL)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	return nil
}

func (client *Client) GenerateResponse(input payload.RestaurantClaim) (*payload.Tokens, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}
	url := fmt.Sprintf("%s%s%s", client.BaseUrl, client.AuthPort, client.GenerateResponseUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	var tokens payload.Tokens
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &tokens, nil
}

func (client *Client) RefreshToken(input payload.RefreshToken) (*payload.Tokens, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}
	url := fmt.Sprintf("%s%s%s", client.BaseUrl, client.AuthPort, client.RefreshTokenUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	var tokens payload.Tokens
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &tokens, nil
}

// working
func (client *Client) ViewRestaurantOrders(input payload.Input) (*[]payload.OrderDetails, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}
	url := fmt.Sprintf("%s%s%s", client.BaseUrl, client.OrderPort, client.RESTAURANT_ORDERS_URL)
	log.Print(url)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	var orders []payload.OrderDetails
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &orders, nil
}

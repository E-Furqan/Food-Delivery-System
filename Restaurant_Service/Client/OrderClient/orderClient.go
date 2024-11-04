package OrderClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

type OrderClient struct {
	environmentVariable.Environment
}

func NewClient(env environmentVariable.Environment) *OrderClient {
	return &OrderClient{
		Environment: env,
	}
}

func (OrderClient *OrderClient) UpdateOrderStatus(input model.OrderDetails, token string) error {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("error marshaling input: %v", err)
	}

	url := fmt.Sprintf("%s%s%s", OrderClient.Environment.BASE_URL, OrderClient.Environment.ORDER_PORT, OrderClient.Environment.UPDATE_ORDER_STATUS_URL)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

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

func (OrderClient *OrderClient) ViewRestaurantOrders(input model.Input, token string) (*[]model.OrderDetails, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}
	url := fmt.Sprintf("%s%s%s", OrderClient.Environment.BASE_URL, OrderClient.Environment.ORDER_PORT, OrderClient.Environment.RESTAURANT_ORDERS_URL)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	var orders []model.OrderDetails
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &orders, nil
}

func (OrderClient *OrderClient) ViewOrdersDetails(input model.OrderDetails, token string) (*model.OrderDetails, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}
	url := fmt.Sprintf("%s%s%s", OrderClient.Environment.BASE_URL, OrderClient.Environment.ORDER_PORT, OrderClient.Environment.VIEW_ORDER_DETAIL_URL)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	var orders model.OrderDetails
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &orders, nil
}

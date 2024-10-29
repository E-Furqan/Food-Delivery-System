package OrderClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
)

type OrderClient struct {
	environmentVariable.Environment
}

func NewClient(env environmentVariable.Environment) *OrderClient {
	return &OrderClient{
		Environment: env,
	}
}
func (orderClient *OrderClient) ProcessOrder(input payload.ProcessOrder) error {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("error marshaling input: %v", err)
	}

	url := fmt.Sprintf("%s%s%s", orderClient.Environment.BASE_URL, orderClient.Environment.ORDER_PORT, orderClient.Environment.PROCESS_ORDER_URL)
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

func (orderClient *OrderClient) ViewUserOrders(input payload.ProcessOrder) (*[]payload.ProcessOrder, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}
	url := fmt.Sprintf("%s%s%s", orderClient.Environment.BASE_URL, orderClient.Environment.ORDER_PORT, orderClient.Environment.User_ORDERS_URL)
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

	var orders []payload.ProcessOrder
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &orders, nil
}
func (orderClient *OrderClient) ViewDriverOrders(input payload.ProcessOrder) (*[]payload.ProcessOrder, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}
	url := fmt.Sprintf("%s%s%s", orderClient.Environment.BASE_URL, orderClient.Environment.ORDER_PORT, orderClient.Environment.DRIVER_ORDERS_URL)
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

	var orders []payload.ProcessOrder
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &orders, nil
}

func (orderClient *OrderClient) ViewOrdersDetails(input payload.ProcessOrder) (*payload.ProcessOrder, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}
	url := fmt.Sprintf("%s%s%s", orderClient.Environment.BASE_URL, orderClient.Environment.ORDER_PORT, orderClient.Environment.VIEW_ORDER_DETAIL_URL)
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

	var orders payload.ProcessOrder
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &orders, nil
}

func (orderClient *OrderClient) ViewOrdersWithoutRider(input payload.ProcessOrder) (*[]payload.ProcessOrder, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}
	url := fmt.Sprintf("%s%s%s", orderClient.Environment.BASE_URL, orderClient.Environment.ORDER_PORT, orderClient.Environment.VIEW_ORDER_WITHOUT_DRIVER_URL)
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

	var orders []payload.ProcessOrder
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &orders, nil
}

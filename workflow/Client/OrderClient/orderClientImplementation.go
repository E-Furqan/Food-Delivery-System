package OrderClient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
)

func (orderClient *OrderClient) UpdateOrderStatus(input model.UpdateOrder, token string) (*model.UpdateOrder, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}

	url, err := utils.CreateUrl(orderClient.OrderClientEnv.BASE_URL,
		orderClient.OrderClientEnv.ORDER_PORT,
		orderClient.OrderClientEnv.UPDATE_ORDER_STATUS_URL)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	req, err := utils.CreateAuthorizedRequest(url, jsonData, "PATCH", token)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	client := utils.CreateHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update order status: received HTTP %d", resp.StatusCode)
	}

	var orders model.UpdateOrder
	limit := int64(1 << 20)
	if err := json.NewDecoder(io.LimitReader(resp.Body, limit)).Decode(&orders); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &orders, nil
}

func (orderClient *OrderClient) AssignDriver(input model.UpdateOrder, token string) error {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("error marshaling input: %v", err)
	}

	url, err := utils.CreateUrl(orderClient.OrderClientEnv.BASE_URL,
		orderClient.OrderClientEnv.ORDER_PORT,
		orderClient.OrderClientEnv.ASSIGN_DRIVER_URL)

	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	req, err := utils.CreateAuthorizedRequest(url, jsonData, "PATCH", token)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	client := utils.CreateHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update order status: received HTTP %d", resp.StatusCode)
	}

	return nil
}

func (orderClient *OrderClient) ViewOrders(input model.UpdateOrder, token string) (*[]model.UpdateOrder, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}

	url, err := utils.CreateUrl(orderClient.OrderClientEnv.BASE_URL,
		orderClient.OrderClientEnv.ORDER_PORT,
		orderClient.OrderClientEnv.VIEW_ORDERS_URL)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	req, err := utils.CreateAuthorizedRequest(url, jsonData, "GET", token)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	client := utils.CreateHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update order status: received HTTP %d", resp.StatusCode)
	}

	var orders []model.UpdateOrder
	limit := int64(1 << 22) //limiting the size of response to 4MB
	if err := json.NewDecoder(io.LimitReader(resp.Body, limit)).Decode(&orders); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &orders, nil
}

func (orderClient *OrderClient) ViewOrdersWithoutDriver(input model.UpdateOrder, token string) (*[]model.UpdateOrder, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}

	url, err := utils.CreateUrl(orderClient.OrderClientEnv.BASE_URL,
		orderClient.OrderClientEnv.ORDER_PORT,
		orderClient.OrderClientEnv.ASSIGN_DRIVER_URL)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	req, err := utils.CreateAuthorizedRequest(url, jsonData, "GET", token)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	client := utils.CreateHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to view orders without driver: received HTTP %d", resp.StatusCode)
	}

	var orders []model.UpdateOrder
	limit := int64(1 << 22) //limiting the size of response to 4MB
	if err := json.NewDecoder(io.LimitReader(resp.Body, limit)).Decode(&orders); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &orders, nil
}

func (orderClient *OrderClient) FetchOrderStatus(orderID uint, token string) (*model.UpdateOrder, error) {

	jsonData, err := json.Marshal(orderID)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}

	url, err := utils.CreateUrl(orderClient.OrderClientEnv.BASE_URL,
		orderClient.OrderClientEnv.ORDER_PORT,
		orderClient.OrderClientEnv.Fetch_OrderStatus_URL)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	req, err := utils.CreateAuthorizedRequest(url, jsonData, "GET", token)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	client := utils.CreateHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update order status: received HTTP %d", resp.StatusCode)
	}

	var orders model.UpdateOrder
	limit := int64(1 << 22) //limiting the size of response to 4MB
	if err := json.NewDecoder(io.LimitReader(resp.Body, limit)).Decode(&orders); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &orders, nil
}

func (orderClient *OrderClient) CreateOrder(order model.CombineOrderItem, token string) (model.UpdateOrder, error) {

	jsonData, err := json.Marshal(order)
	if err != nil {
		return model.UpdateOrder{}, fmt.Errorf("error marshaling input: %v", err)
	}

	url, err := utils.CreateUrl(orderClient.OrderClientEnv.BASE_URL,
		orderClient.OrderClientEnv.ORDER_PORT,
		orderClient.OrderClientEnv.Create_Order_URL)
	if err != nil {
		log.Printf("BASE_URL: %s", orderClient.OrderClientEnv.BASE_URL)
		log.Printf("ORDER_PORT: %s", orderClient.OrderClientEnv.ORDER_PORT)
		log.Printf("Create_Order_URL: %s", orderClient.OrderClientEnv.Create_Order_URL)
		return model.UpdateOrder{}, fmt.Errorf("error: %v", err)
	}

	req, err := utils.CreateAuthorizedRequest(url, jsonData, "POST", token)
	if err != nil {
		return model.UpdateOrder{}, fmt.Errorf("error: %v", err)
	}

	client := utils.CreateHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return model.UpdateOrder{}, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.UpdateOrder{}, fmt.Errorf("failed to update order status: received HTTP %d", resp.StatusCode)
	}

	var orders model.UpdateOrder
	limit := int64(1 << 22) //limiting the size of response to 4MB
	if err := json.NewDecoder(io.LimitReader(resp.Body, limit)).Decode(&orders); err != nil {
		return model.UpdateOrder{}, fmt.Errorf("error decoding response: %v", err)
	}

	return orders, nil
}

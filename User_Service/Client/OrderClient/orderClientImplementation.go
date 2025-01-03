package OrderClient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

func (orderClient *OrderClient) UpdateOrderStatus(input model.UpdateOrder, c *gin.Context) (*model.UpdateOrder, error) {

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

	req, err := utils.CreateAuthorizedRequest(url, jsonData, c, "PATCH")
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

func (orderClient *OrderClient) AssignDriver(input model.UpdateOrder, c *gin.Context) error {

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

	req, err := utils.CreateAuthorizedRequest(url, jsonData, c, "PATCH")
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

func (orderClient *OrderClient) ViewOrders(input model.UpdateOrder, c *gin.Context) (*[]model.UpdateOrder, error) {

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

	req, err := utils.CreateAuthorizedRequest(url, jsonData, c, "GET")
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

func (orderClient *OrderClient) ViewOrdersWithoutDriver(input model.UpdateOrder, c *gin.Context) (*[]model.UpdateOrder, error) {

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

	req, err := utils.CreateAuthorizedRequest(url, jsonData, c, "GET")
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

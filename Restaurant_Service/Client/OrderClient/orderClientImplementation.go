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

func (OrderClient *OrderClient) UpdateOrderStatus(input model.OrderDetails, c *gin.Context) error {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("error marshaling input: %v", err)
	}

	url, err := utils.CreateUrl(OrderClient.OrderClientEnv.BASE_URL, OrderClient.OrderClientEnv.ORDER_PORT, OrderClient.OrderClientEnv.UPDATE_ORDER_STATUS_URL)
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
		return fmt.Errorf("failed to update orders status: received HTTP %d", resp.StatusCode)
	}
	return nil
}

func (OrderClient *OrderClient) ViewRestaurantOrders(input model.Input, c *gin.Context) (*[]model.OrderDetails, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}

	url, err := utils.CreateUrl(OrderClient.OrderClientEnv.BASE_URL, OrderClient.OrderClientEnv.ORDER_PORT, OrderClient.OrderClientEnv.RESTAURANT_ORDERS_URL)
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
		return nil, fmt.Errorf("failed to view restaurant orders: received HTTP %d", resp.StatusCode)
	}

	var orders []model.OrderDetails
	limit := int64(1 << 22) //limiting the size of response to 4MB
	if err := json.NewDecoder(io.LimitReader(resp.Body, limit)).Decode(&orders); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &orders, nil
}

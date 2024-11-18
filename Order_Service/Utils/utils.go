package utils

import (
	"fmt"
	"os"
	"strings"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/gin-gonic/gin"
)

func GenerateResponse(httpStatusCode int, c *gin.Context, title1 string, message1 string, title2 string, input interface{}) {
	response := gin.H{
		title1: message1,
	}

	if title2 != "" && input != nil {
		response[title2] = input
	}

	c.JSON(httpStatusCode, response)
}

func GetEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func CreateInvoice(order model.Order, orderItems []model.OrderItem, items []model.Items) gin.H {
	invoiceItems := []gin.H{}
	totalBill := order.TotalBill

	for _, orderItem := range orderItems {
		for _, item := range items {
			if item.ItemId == orderItem.ItemId {
				invoiceItems = append(invoiceItems, gin.H{
					"item_id":    item.ItemId,
					"name":       item.ItemName,
					"quantity":   orderItem.Quantity,
					"unit_price": item.ItemPrice,
					"total":      float64(orderItem.Quantity) * item.ItemPrice,
				})
			}
		}

	}

	return gin.H{
		"order_id":      order.OrderID,
		"user_id":       order.UserId,
		"restaurant_id": order.RestaurantID,
		"order_status":  order.OrderStatus,
		"total_bill":    totalBill,
		"items":         invoiceItems,
	}
}
func CreateOrderObj(order model.CombineOrderItem, bill float64) model.Order {
	return model.Order{
		OrderStatus:  "order placed",
		UserId:       order.UserId,
		RestaurantID: order.RestaurantId,
		TotalBill:    bill,
	}
}

func CalculateBill(CombineOrderItem model.CombineOrderItem, items []model.Items) (float64, error) {
	totalBill := 0.0

	for _, orderedItem := range CombineOrderItem.Items {
		var ItemPrice float64
		ItemFound := false

		for _, item := range items {
			if item.ItemId == orderedItem.ItemId {
				ItemPrice = item.ItemPrice
				ItemFound = true
				break
			}
		}

		if !ItemFound {
			continue
		}

		totalBill += ItemPrice * float64(orderedItem.Quantity)
	}
	if totalBill == 0 {
		return totalBill, fmt.Errorf("items are not of this restaurant")
	}
	return totalBill, nil
}

func VerifyRole(c *gin.Context) (string, error) {
	activeRole, exists := c.Get("activeRole")
	if !exists {
		return "", fmt.Errorf("userId role does not exist")
	}

	activeRoleStr, ok := activeRole.(string)
	if !ok {
		return "", fmt.Errorf("activeRole is not a string")
	}
	return activeRoleStr, nil
}

func IsCustomerOrAdminRole(activeRoleStr string) bool {
	if strings.ToLower(activeRoleStr) == "customer" || strings.ToLower(activeRoleStr) == "admin" {
		return true
	} else {
		return false
	}
}

func IsRestaurantOrAdminRole(activeRoleStr string) bool {
	if strings.ToLower(activeRoleStr) == "restaurant" || strings.ToLower(activeRoleStr) == "admin" {
		return true
	} else {
		return false
	}
}

func IsDriverOrAdminRole(activeRoleStr string) bool {
	if strings.ToLower(activeRoleStr) == "delivery driver" || strings.ToLower(activeRoleStr) == "admin" {
		return true
	} else {
		return false
	}
}

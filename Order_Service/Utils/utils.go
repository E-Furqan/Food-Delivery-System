package utils

import (
	"os"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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

type Claims struct {
	UserId      uint   `json:"user_id"`
	ActiveRole  string `json:"activeRole"`
	ClaimId     uint   `json:"claim_id"`
	ServiceType string `json:"service_type"`
	jwt.StandardClaims
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

	return totalBill, nil
}

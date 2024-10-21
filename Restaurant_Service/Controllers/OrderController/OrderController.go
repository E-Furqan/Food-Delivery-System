package OrderController

import (
	"log"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	Repo *database.Repository
}

func NewController(repo *database.Repository) *OrderController {
	return &OrderController{Repo: repo}
}

func (OrderController *OrderController) ProcessOrder(c *gin.Context) {
	var order payload.ProcessOrder

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, "Error while binding order status")
		return
	}

	orderTransitions := payload.GetOrderTransitions()

	if order.OrderStatus == "order placed" {
		var restaurant model.Restaurant
		err := OrderController.Repo.GetRestaurant("restaurant_id", order.RestaurantId, &restaurant)
		if err != nil {
			order.OrderStatus = "Cancelled"
			c.JSON(http.StatusNotFound, "Restaurant not found")
			c.JSON(http.StatusNotFound, order)
			return
		}

		if restaurant.RestaurantStatus == "closed" || restaurant.RestaurantStatus == "Closed" {
			order.OrderStatus = "Cancelled"
			c.JSON(http.StatusBadRequest, "Restaurant is closed")
			c.JSON(http.StatusBadRequest, order)
			return
		}
	}

	if newStatus, exists := orderTransitions[order.OrderStatus]; exists {
		order.OrderStatus = newStatus
	}

	c.JSON(http.StatusOK, order)
}

func (OrderController *OrderController) CancelOrder(c *gin.Context) {
	email, exists := c.Get("Email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant not authenticated"})
		return
	}
	email, ok := email.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Email address"})
		return
	}

	var Restaurant model.Restaurant
	err := OrderController.Repo.GetRestaurant("restaurant_email", email, &Restaurant)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var input payload.ProcessOrder

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Binding input data failed",
			"Error":   err,
		})
		return
	}

	if input.RestaurantId != Restaurant.RestaurantId {
		log.Print(Restaurant.RestaurantId)
		log.Print(input.RestaurantId)
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "You are not authorized to cancel this order as it belongs to a different restaurant",
		})
		return
	}

	input.OrderStatus = "Cancelled"
	c.JSON(http.StatusOK, gin.H{
		"Message":       "Order cancelled successfully",
		"Order details": input,
	})

}

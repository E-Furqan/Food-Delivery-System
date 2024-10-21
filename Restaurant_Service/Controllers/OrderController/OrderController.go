package OrderController

import (
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

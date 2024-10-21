package OrderController

import (
	"net/http"

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

	if order.OrderStatus == "Cancelled" {
		c.JSON(http.StatusBadRequest, "Order is cancelled")
		c.JSON(http.StatusBadRequest, order)
		return
	}

	if order.OrderStatus == "Waiting For Delivery Driver" {
		driver, err := OrderController.Repo.GetDeliveryDrivers()
		if err != nil {
			c.JSON(http.StatusNotFound, "Delivery Driver Not Found")
			return
		}
		if newStatus, exists := orderTransitions[order.OrderStatus]; exists {
			order.OrderStatus = newStatus
		}

		c.JSON(http.StatusOK, gin.H{
			"Order":          order,
			"Deliver Driver": driver.UserId,
		})

	}

	if newStatus, exists := orderTransitions[order.OrderStatus]; exists {
		order.OrderStatus = newStatus
	}

	c.JSON(http.StatusOK, gin.H{
		"Order": order,
	})
}

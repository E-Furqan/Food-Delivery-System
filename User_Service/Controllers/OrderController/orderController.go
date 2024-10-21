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

	if order.OrderStatus == "Cancelled" {
		c.JSON(http.StatusBadRequest, "Order is cancelled")
		c.JSON(http.StatusBadRequest, order)
		return
	}

	if order.OrderStatus == "Waiting For Delivery Driver" {
		var driver model.User
		err := OrderController.Repo.GetDeliveryDrivers(&driver)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Message": "Delivery driver not found",
				"Error":   err.Error(),
			})
			return
		}

		if newStatus, exists := orderTransitions[order.OrderStatus]; exists {
			order.OrderStatus = newStatus
		}
		order.DeliverDriverID = driver.UserId

		c.JSON(http.StatusOK, gin.H{
			"Order":          order,
			"Deliver Driver": driver.UserId,
		})
		return

	} else if order.OrderStatus == "Delivered" {
		var driver model.User
		log.Print(order.DeliverDriverID)
		err := OrderController.Repo.GetUser("user_id", order.DeliverDriverID, &driver)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error":   err.Error(),
				"Message": "Delivery rider not found",
			})
		}

		driver.RoleStatus = "available"
		OrderController.Repo.UpdateRoleStatus(&driver)
	}

	if newStatus, exists := orderTransitions[order.OrderStatus]; exists {
		log.Print(newStatus)
		order.OrderStatus = newStatus
	}

	c.JSON(http.StatusOK, gin.H{
		"Order": order,
	})
}

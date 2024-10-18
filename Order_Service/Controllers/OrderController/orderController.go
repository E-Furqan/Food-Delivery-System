package OrderControllers

import (
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
)

// Controller struct that holds a reference to the repository
type OrderController struct {
	Repo *database.Repository
}

// NewController initializes the controller with the repository dependency
func NewController(repo *database.Repository) *OrderController {
	return &OrderController{Repo: repo}
}

func (orderCtrl *OrderController) CheckOut(c *gin.Context) {
	var OrderDetails model.OrderItem
	if err := c.ShouldBindJSON(&OrderDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// total_bill := OrderDetails.ItemPrice * OrderDetails.Quantity

}

func (orderCtrl *OrderController) UpdateOrderStatus(c *gin.Context) {

	var OrderStatus payload.Order
	if err := c.ShouldBindJSON(&OrderStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var order model.Order

	err := orderCtrl.Repo.GetOrder(&order, int(OrderStatus.OrderID))
	if err != nil {
		c.JSON(http.StatusNotFound, "Order not found")
		return
	}

	if err := orderCtrl.Repo.Update(&order, OrderStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (orderCtrl *OrderController) GetOrdersOfUser(c *gin.Context) {

	var OrderPayLoad payload.Order
	if err := c.ShouldBindJSON(&OrderPayLoad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var order []model.Order
	if err := orderCtrl.Repo.GetOrders(&order, int(OrderPayLoad.UserId)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)

}

func (orderCtrl *OrderController) PutOrder(c *gin.Context) {

	var CombineOrderItem payload.CombineOrderItem
	if err := c.ShouldBindJSON(&CombineOrderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var order model.Order

	order.OrderID = CombineOrderItem.Order.OrderID
	order.RestaurantID = CombineOrderItem.Order.RestaurantID
	order.TotalBill = 0
	order.OrderStatus = CombineOrderItem.Order.OrderStatus

	err := orderCtrl.Repo.PutOrder(&order, &CombineOrderItem)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while creating order",
			"Error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order created successfully",
		"order":   order,
	})
}

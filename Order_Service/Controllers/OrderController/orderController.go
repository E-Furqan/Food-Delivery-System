package OrderControllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/E-Furqan/Food-Delivery-System/Client/RestaurantClient"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	Repo      *database.Repository
	ResClient *RestaurantClient.RestaurantClient
}

func NewController(repo *database.Repository, ResClient *RestaurantClient.RestaurantClient) *OrderController {
	return &OrderController{
		Repo:      repo,
		ResClient: ResClient,
	}
}

func (orderCtrl *OrderController) UpdateOrderStatus(c *gin.Context) {
	var request model.OrderStatusUpdateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order model.Order

	err := orderCtrl.Repo.GetOrder(&order, request.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if strings.ToLower(request.OrderStatus) == "cancelled" {
		order.OrderStatus = request.OrderStatus
		if err := orderCtrl.Repo.Update(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, order)
		return
	}

	orderTransitions := model.GetOrderTransitions()

	switch strings.ToLower(request.Role) {
	case "user":
		if newStatus, exists := orderTransitions["user"][order.OrderStatus]; exists {
			order.OrderStatus = newStatus
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not allowed to update the order status at this point"})
			return
		}

	case "restaurant":
		if newStatus, exists := orderTransitions["restaurant"][order.OrderStatus]; exists {
			order.OrderStatus = newStatus
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Restaurant is not allowed to update the order status at this point"})
			return
		}

	case "delivery driver":
		if newStatus, exists := orderTransitions["delivery driver"][order.OrderStatus]; exists {
			order.OrderStatus = newStatus
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Delivery driver is not allowed to update the order status at this point"})
			return
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	if err := orderCtrl.Repo.Update(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (orderCtrl *OrderController) AssignDeliveryDriver(c *gin.Context) {
	var request model.AssignDeliveryDriver
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.ToLower(request.Role) != "delivery driver" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	var order model.Order

	err := orderCtrl.Repo.GetOrder(&order, request.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.DeliveryDriverID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order already have a driver"})
		return
	}

	order.DeliveryDriverID = request.DeliveryDriverID
	if err := orderCtrl.Repo.Update(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (orderCtrl *OrderController) GetOrders(c *gin.Context, UserType string) {

	var OrderNFilter model.CombineOrderFilter
	if err := c.ShouldBindJSON(&OrderNFilter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order []model.Order
	var err error

	if UserType == "user" {
		err = orderCtrl.Repo.GetOrders(&order, OrderNFilter.UserId, OrderNFilter.Filter.ColumnName, OrderNFilter.Filter.OrderDirection, "user_id")
	} else if UserType == "restaurant" {
		err = orderCtrl.Repo.GetOrders(&order, OrderNFilter.RestaurantId, OrderNFilter.Filter.ColumnName, OrderNFilter.Filter.OrderDirection, "restaurant_id")

	} else if UserType == "delivery driver" {
		err = orderCtrl.Repo.GetOrders(&order, OrderNFilter.RestaurantId, OrderNFilter.Filter.ColumnName, OrderNFilter.Filter.OrderDirection, "delivery_driver")
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (orderCtrl *OrderController) GetOrdersOfUser(c *gin.Context) {
	orderCtrl.GetOrders(c, "user")
}

func (orderCtrl *OrderController) GetOrdersOfRestaurant(c *gin.Context) {
	orderCtrl.GetOrders(c, "restaurant")
}
func (orderCtrl *OrderController) GetOrdersOfDeliveryDriver(c *gin.Context) {
	orderCtrl.GetOrders(c, "delivery driver")
}

func (orderCtrl *OrderController) PlaceOrder(c *gin.Context) {
	ServiceType, exists := c.Get("ServiceType")
	if !exists {
		c.JSON(http.StatusBadRequest, "userId id does not exist")
		return
	}
	if ServiceType != "User" {
		c.JSON(http.StatusBadRequest, "Only user can place order")
		return
	}

	var CombineOrderItem model.CombineOrderItem
	if err := c.ShouldBindJSON(&CombineOrderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var GetItem model.GetItems
	GetItem.RestaurantId = CombineOrderItem.RestaurantId
	GetItem.ColumnName = "restaurant_id"
	GetItem.OrderType = "asc"

	var items []model.Items
	items, err := orderCtrl.ResClient.GetItems(GetItem)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Error while getting items from the restaurant", "Error", err.Error())
		return
	}

	if len(items) == 0 {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "No items found in the restaurant", "", nil)
		return
	}

	totalBill, err := utils.CalculateBill(CombineOrderItem, items)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	order := utils.CreateOrderObj(CombineOrderItem, totalBill)
	err = orderCtrl.Repo.PlaceOrder(&order, &CombineOrderItem)

	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Error while creating order", "Error", err.Error())
		return
	}
	utils.GenerateResponse(http.StatusOK, c, "Message", fmt.Sprintf("Order created successfully with total bill: %v", totalBill), "", nil)
}

func (orderCtrl *OrderController) ViewOrderDetails(c *gin.Context) {
	var orderId model.ID

	if err := c.ShouldBindJSON(&orderId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var order model.Order
	err := orderCtrl.Repo.GetOrder(&order, orderId.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, "Order not found")
		return
	}

	c.JSON(http.StatusOK, order)
}

func (orderCtrl *OrderController) ViewOrdersWithoutRider(c *gin.Context) {

	var order []model.Order
	err := orderCtrl.Repo.GetOrderWithoutRider(&order)
	if err != nil {
		c.JSON(http.StatusNotFound, "Order not found")
		return
	}

	c.JSON(http.StatusOK, order)

}

func (orderCtrl *OrderController) GenerateInvoice(c *gin.Context) {

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusBadRequest, "userId id does not exist")
		return
	}

	var orderId model.ID
	if err := c.ShouldBindJSON(&orderId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order model.Order
	if err := orderCtrl.Repo.GetOrder(&order, orderId.OrderID); err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Message", "Order not found", "Error", err.Error())
		return
	}

	if userId != order.UserId {
		utils.GenerateResponse(http.StatusNotFound, c, "Error", "You are not allowed to generate invoice of this order", "", nil)
		return
	}

	var GetItem model.GetItems
	GetItem.RestaurantId = order.RestaurantID
	GetItem.ColumnName = "restaurant_id"
	GetItem.OrderType = "asc"

	var items []model.Items
	items, err := orderCtrl.ResClient.GetItems(GetItem)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Error while getting items from the restaurant", "Error", err.Error())
		return
	}

	var orderItems []model.OrderItem
	if err := orderCtrl.Repo.GetOrderItems(&orderItems, orderId.OrderID); err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Message", "Error retrieving order items", "Error", err.Error())
		return
	}

	invoice := utils.CreateInvoice(order, orderItems, items)

	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}

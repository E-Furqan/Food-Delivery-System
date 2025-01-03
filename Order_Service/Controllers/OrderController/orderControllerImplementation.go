package OrderControllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

// UpdateOrderStatus godoc
// @Summary Update the status of an order
// @Description Updates the status of an order based on the role of the user (customer, restaurant, or delivery driver)
// @Tags Order Service
// @Accept  json
// @Produce  json
// @Param orderStatusUpdateRequest body model.OrderStatusUpdateRequest true "Order status update request"
// @Success 200 {object} model.Order "Updated order details"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Router /order/update/status [patch]
func (orderCtrl *OrderController) UpdateOrderStatus(c *gin.Context) {
	Id, exists := c.Get("ClaimId")
	if !exists {
		c.JSON(http.StatusBadRequest, "userId id does not exist")
		return
	}

	activeRoleStr, err := utils.VerifyRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	var request model.OrderStatusUpdateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order model.Order

	err = orderCtrl.Repo.GetOrder(&order, request.OrderID)
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

	switch strings.ToLower(activeRoleStr) {
	case "customer":

		if order.UserId != Id {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have the permission to update order status"})
			return
		}

		if newStatus, exists := orderTransitions["user"][order.OrderStatus]; exists {
			order.OrderStatus = newStatus
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not allowed to update the order status at this point"})
			return
		}

	case "restaurant":

		if order.RestaurantID != Id {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have the permission to update order status"})
			return
		}

		if newStatus, exists := orderTransitions["restaurant"][order.OrderStatus]; exists {
			order.OrderStatus = newStatus
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant is not allowed to update the order status at this point"})
			return
		}

	case "delivery driver":

		if order.DeliveryDriverID != Id {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have the permission to update order status"})
			return
		}

		if newStatus, exists := orderTransitions["delivery driver"][order.OrderStatus]; exists {
			order.OrderStatus = newStatus
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Delivery driver is not allowed to update the order status at this point"})
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

// AssignDeliveryDriver godoc
// @Summary Assign a delivery driver to an order
// @Description Assigns a delivery driver to an order if the order doesn't already have a driver, and the role of the user is "delivery driver"
// @Tags Order Service
// @Accept  json
// @Produce  json
// @Param assignDeliveryDriverRequest body model.AssignDeliveryDriver true "Assign delivery driver request"
// @Success 200 {object} model.Order "Assigned delivery driver to the order"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Router /order/assign/diver [patch]
func (orderCtrl *OrderController) AssignDeliveryDriver(c *gin.Context) {
	Id, exists := c.Get("ClaimId")
	if !exists {
		c.JSON(http.StatusBadRequest, "userId id does not exist")
		return
	}
	IDint := Id.(uint)

	activeRoleStr, err := utils.VerifyRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	var request model.AssignDeliveryDriver
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.ToLower(activeRoleStr) != "delivery driver" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	var order model.Order

	err = orderCtrl.Repo.GetOrder(&order, request.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.DeliveryDriverID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order already have a driver"})
		return
	}
	order.DeliveryDriverID = IDint
	if err := orderCtrl.Repo.Update(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrders retrieves orders based on the user type.
//
// @Summary Get orders based on user type (user, restaurant, or delivery driver)
// @Description Retrieves orders filtered by user type, sorted by specified column and order.
// @Tags Order Service
// @Param ID header uint true "User ID from Claims"
// @Param UserType path string true "Type of user: user, restaurant, or delivery driver"
// @Param Filter body model.Filter true "Sorting details"
// @Success 200 {array} model.Order "List of Orders"
// @Failure 400 {object} model.ErrorResponse "Error occurred"
// @Router /order/view/orders [get]
// @Security ApiKeyAuth
func (orderCtrl *OrderController) GetOrders(c *gin.Context) {

	activeRoleStr, err := utils.VerifyRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	Id, exists := c.Get("ID")
	if !exists {
		c.JSON(http.StatusBadRequest, "userId id does not exist")
		return
	}

	IdValue := Id.(uint)
	var Filter model.Filter
	if err := c.ShouldBindJSON(&Filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order []model.Order

	if utils.IsCustomerOrAdminRole(activeRoleStr) {
		err = orderCtrl.Repo.GetOrders(&order, IdValue, Filter.ColumnName, Filter.SortOrder, "user_id")

	} else if utils.IsRestaurantOrAdminRole(activeRoleStr) {
		err = orderCtrl.Repo.GetOrders(&order, IdValue, Filter.ColumnName, Filter.SortOrder, "restaurant_id")

	} else if utils.IsDriverOrAdminRole(activeRoleStr) {
		err = orderCtrl.Repo.GetOrders(&order, IdValue, Filter.ColumnName, Filter.SortOrder, "delivery_driver")
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// PlaceOrder godoc
// @Summary Place an order as a customer
// @Description Allows a customer to place an order, including selecting items from a restaurant and calculating the total bill
// @Tags Order Service
// @Accept  json
// @Produce  json
// @Param placeOrderRequest body model.CombineOrderItem true "Place order request"
// @Success 200 {object} map[string]interface{} "Order placed successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /order/place/order [post]
func (orderCtrl *OrderController) PlaceOrder(c *gin.Context) {

	userRoleStr, err := utils.VerifyRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	if strings.ToLower(userRoleStr) != "customer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only customer can place order"})
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
	items, err = orderCtrl.ResClient.GetItems(GetItem)
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
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	order := utils.CreateOrderObj(CombineOrderItem, totalBill)
	err = orderCtrl.Repo.PlaceOrder(&order, &CombineOrderItem)

	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Error while creating order", "Error", err.Error())
		return
	}
	utils.GenerateResponse(http.StatusOK, c, "Message", fmt.Sprintf("Order created successfully with order id %v and total bill: %v", order.OrderID, totalBill), "", nil)
}

// ViewOrderDetails godoc
// @Summary View details of a specific order
// @Description Retrieves detailed information about an order by order ID
// @Tags Order Service
// @Accept  json
// @Produce  json
// @Param orderId body model.ID true "Order ID JSON"
// @Success 200 {object} model.Order "Order details"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} string "Order not found"
// @Router /order/view/order [get]
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

// ViewOrdersWithoutRider godoc
// @Summary Get orders without assigned delivery driver
// @Description Retrieves orders that have not been assigned a delivery driver. Only accessible to users with roles "delivery driver" or "admin".
// @Tags Order Service
// @Produce json
// @Success 200 {array} model.Order "List of orders without assigned delivery driver"
// @Failure 400 {object} model.ErrorResponse "Bad request, role ID missing or unauthorized access"
// @Failure 404 {object} model.ErrorResponse "Orders not found"
// @Router /view/without/driver/orders [get]
// @Security ApiKeyAuth
func (orderCtrl *OrderController) ViewOrdersWithoutRider(c *gin.Context) {

	activeRoleStr, err := utils.VerifyRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	if strings.ToLower(activeRoleStr) != "delivery driver" && strings.ToLower(activeRoleStr) != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only delivery driver and admin can view the orders"})
		return
	}

	var order []model.Order
	err = orderCtrl.Repo.GetOrderWithoutRider(&order)
	if err != nil {
		c.JSON(http.StatusNotFound, "Order not found")
		return
	}

	c.JSON(http.StatusOK, order)
}

// GenerateInvoice godoc
// @Summary Generate an invoice for a specific order
// @Description Generates an invoice for the order based on order ID and user validation
// @Tags Order Service
// @Accept  json
// @Produce  json
// @Param orderId body model.ID true "Order ID JSON"
// @Success 200 {object} map[string]interface{} "invoice"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /order/generate/invoice [get]
func (orderCtrl *OrderController) GenerateInvoice(c *gin.Context) {

	userId, exists := c.Get("ClaimId")
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
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "You are not allowed to generate invoice of this order", "", nil)
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
	log.Print(orderItems)
	invoice := utils.CreateInvoice(order, orderItems, items)

	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}

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

	activeRoleStr, err := utils.FetchRoleFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
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

		if order.DeliveryDriver != Id {
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

	activeRoleStr, err := utils.FetchRoleFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
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

	if order.DeliveryDriver != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order already have a driver"})
		return
	}
	order.DeliveryDriver = IDint
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

	activeRoleStr, err := utils.FetchRoleFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
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

	userRoleStr, err := utils.FetchRoleFromClaims(c)
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

	activeRoleStr, err := utils.FetchRoleFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
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

// FetchAverageOrderValue godoc
// @Summary Fetch average order value for a user, restaurant, or time period
// @Description Retrieves the average order value based on user ID, restaurant ID, or a time period filter
// @Tags Order Service
// @Accept  json
// @Produce  json
// @Param input body model.AverageOrderValue true "Input parameters for filtering"
// @Success 200 {object} map[string]interface{} "Average order value"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /order/average/order_value [get]
// @Security ApiKeyAuth
func (orderCtrl *OrderController) FetchAverageOrderValue(c *gin.Context) {

	var input model.AverageOrderValue
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activeRoleStr, err := utils.FetchRoleFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ID, err := utils.FetchIDFromClaim(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, "ClaimId is not a valid uint")
		return
	}

	if input.FilterType == "user" && utils.IsCustomerOrAdminRole(activeRoleStr) {
		result, err := orderCtrl.Repo.FetchAverageOrderValueOfUser(ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
		return

	} else if input.FilterType == "restaurant" && utils.IsRestaurantOrAdminRole(activeRoleStr) {
		result, err := orderCtrl.Repo.FetchAverageOrderValueOfRestaurant(ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
		return

	} else if input.FilterType == "time" && utils.IsAdminRole(activeRoleStr) {
		result, err := orderCtrl.Repo.FetchAverageOrderValueBetweenTime(input.StartTime, input.EndTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}
}

// FetchCompletedDeliversRider godoc
// @Summary Fetch completed deliveries by riders
// @Description Retrieves the number of completed deliveries for each rider. This endpoint is restricted to admins only.
// @Tags Order Service
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.CompletedDelivers "List of completed deliveries by riders"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /order/completed/delivers/rider [get]
// @Summary Fetch completed deliveries by riders
// @Description Retrieves the number of completed deliveries for each rider
// @Tags Order Service
// @Accept  json
// @Produce  json
// @Success 200 {array} model.CompletedDelivers "List of completed deliveries by riders"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /order/completed/delivers/rider [get]
// @Security ApiKeyAuth
func (orderCtrl *OrderController) FetchCompletedDeliversRider(c *gin.Context) {

	activeRoleStr, err := utils.FetchRoleFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	isAdmin := utils.IsAdminRole(activeRoleStr)
	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can get completed delivers by riders analysis"})
		return
	}

	result, err := orderCtrl.Repo.FetchCompletedDeliversOfRider()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// FetchCancelOrdersDetails godoc
// @Summary Fetch cancelled orders details
// @Description Retrieves the details of cancelled orders including item information
// @Tags Order Service
// @Accept json
// @Produce json
// @Param pageNumber body model.PageNumber true "Pagination information"
// @Success 200 {array} model.OrderDetails "List of cancelled orders with item details"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid page number or limit"
// @Failure 401 {object} map[string]interface{} "Unauthorized: User not authorized to access the resource"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /order/cancelled/orders [get]
// @Security ApiKeyAuth
func (orderCtrl *OrderController) FetchCancelOrdersDetails(c *gin.Context) {

	activeRoleStr, err := utils.FetchRoleFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	isAdmin := utils.IsAdminRole(activeRoleStr)
	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can get completed delivers by riders analysis"})
		return
	}

	var PageNumber model.PageNumber
	if err := c.ShouldBindJSON(&PageNumber); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	offset := (PageNumber.PageNumber - 1) * PageNumber.Limit
	result, err := orderCtrl.Repo.FetchCancelledOrdersWithItemDetails(PageNumber.Limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// FetchCustomerOrdersDetails godoc
// @Summary Fetch orders details for a specific customer
// @Description Retrieves the details of orders placed by a specific customer, including item information
// @Tags Order Service
// @Accept json
// @Produce json
// @Param pageNumber body model.PageNumber true "Pagination information"
// @Success 200 {array} model.OrderDetails "List of orders placed by the customer with item details"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid page number or limit"
// @Failure 401 {object} map[string]interface{} "Unauthorized: User not authorized to access the resource"
// @Failure 404 {object} map[string]interface{} "Not Found: Customer not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /order/customer/orders [get]
// @Security ApiKeyAuth
func (orderCtrl *OrderController) FetchCustomerOrdersDetails(c *gin.Context) {

	activeRoleStr, err := utils.FetchRoleFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	isAdmin := utils.IsCustomerOrAdminRole(activeRoleStr)
	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can get completed delivers by riders analysis"})
		return
	}

	ID, err := utils.FetchIDFromClaim(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, "ClaimId is not a valid uint")
		return
	}

	var PageNumber model.PageNumber
	if err := c.ShouldBindJSON(&PageNumber); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	offset := (PageNumber.PageNumber - 1) * PageNumber.Limit
	result, err := orderCtrl.Repo.FetchUserOrdersWithItemDetails(int(ID), PageNumber.Limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (orderCtrl *OrderController) FetchTopPurchasedItems(c *gin.Context) {

	result, err := orderCtrl.Repo.FetchTopPurchasedItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (orderCtrl *OrderController) FetchCompletedOrdersCountByRestaurant(c *gin.Context) {

	activeRoleStr, err := utils.FetchRoleFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	isAdminOrRestaurant := utils.IsRestaurantOrAdminRole(activeRoleStr)
	if !isAdminOrRestaurant {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can get completed delivers by riders analysis"})
		return
	}

	var TimeRange model.TimeRange
	if err := c.ShouldBindJSON(&TimeRange); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	result, err := orderCtrl.Repo.FetchCompletedOrdersCountByRestaurant(TimeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (orderCtrl *OrderController) FetchOrderStatusFrequencies(c *gin.Context) {
	// activeRoleStr, err := utils.FetchRoleFromClaims(c)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	return
	// }

	// isAdmin := utils.IsAdminRole(activeRoleStr)
	// if !isAdmin {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can get completed delivers by riders analysis"})
	// 	return
	// }

	result, err := orderCtrl.Repo.FetchOrderStatusFrequencies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

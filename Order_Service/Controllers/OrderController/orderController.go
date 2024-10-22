package OrderControllers

import (
	"fmt"
	"net/http"

	ClientPackage "github.com/E-Furqan/Food-Delivery-System/Client"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

// Controller struct that holds a reference to the repository
type OrderController struct {
	Repo   *database.Repository
	Client *ClientPackage.Client
}

// NewController initializes the controller with the repository dependency
func NewController(repo *database.Repository, client *ClientPackage.Client) *OrderController {
	return &OrderController{
		Repo:   repo,
		Client: client,
	}
}

// func (orderCtrl *OrderController) CheckOut(c *gin.Context) {
// 	var inputOrderId payload.Order
// 	if err := c.ShouldBindJSON(&inputOrderId); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"Error While binding ": err.Error()})
// 		return
// 	}

// 	var orderDetail model.Order

// 	if err := orderCtrl.Repo.GetOrder(&orderDetail, int(inputOrderId.OrderID)); err == nil {
// 		c.JSON(http.StatusNotFound, "Order Not found")
// 		return
// 	}

// 	c.JSON(http.StatusOK, orderDetail)
// }

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

func (orderCtrl *OrderController) GetOrders(c *gin.Context, isUser bool) {

	var OrderNFilter payload.CombineOrderFilter
	if err := c.ShouldBindJSON(&OrderNFilter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order []model.Order
	var err error

	if isUser {
		err = orderCtrl.Repo.GetOrders(&order, int(OrderNFilter.Order.UserId), OrderNFilter.Filter.ColumnName, OrderNFilter.Filter.OrderDirection)
	} else {
		err = orderCtrl.Repo.GetOrders(&order, int(OrderNFilter.Order.RestaurantID), OrderNFilter.Filter.ColumnName, OrderNFilter.Filter.OrderDirection)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (orderCtrl *OrderController) GetOrdersOfUser(c *gin.Context) {
	orderCtrl.GetOrders(c, true)
}

func (orderCtrl *OrderController) GetOrdersOfRestaurant(c *gin.Context) {
	orderCtrl.GetOrders(c, false)
}

func (orderCtrl *OrderController) PlaceOrder(c *gin.Context) {

	var CombineOrderItem payload.CombineOrderItem
	if err := c.ShouldBindJSON(&CombineOrderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	items, err := orderCtrl.Client.GetItems(CombineOrderItem.Order.RestaurantID)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Error while getting items from the restaurant", "Error", err.Error())
		return
	}

	if len(items) == 0 {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "No items found in the restaurant", "", nil)
		return
	}

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
			utils.GenerateResponse(http.StatusBadRequest, c, "Message", fmt.Sprintf("Item with ID %d not found", orderedItem.ItemId), "", nil)
			return
		}

		totalBill += ItemPrice * float64(orderedItem.Quantity)
	}

	var order model.Order
	order.UserId = CombineOrderItem.Order.UserId
	order.RestaurantID = CombineOrderItem.Order.RestaurantID
	order.TotalBill = totalBill
	order.OrderStatus = "order placed"

	err = orderCtrl.Repo.PlaceOrder(&order, &CombineOrderItem)

	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Error while creating order", "Error", err.Error())
		return
	}
	utils.GenerateResponse(http.StatusOK, c, "Message", "Order created successfully", "", nil)
	var processOrder payload.ProcessOrder
	processOrder.OrderStatus = order.OrderStatus
	processOrder.RestaurantId = order.RestaurantID

	err = orderCtrl.Client.ProcessOrder(processOrder)

	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Error sending request to restaurant service", "Error", err.Error())
		return
	}

	utils.GenerateResponse(http.StatusOK, c, "Message", "Order Accepted by the restaurant successfully", "", nil)

}

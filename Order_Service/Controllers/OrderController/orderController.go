package OrderControllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	RestaurantClient "github.com/E-Furqan/Food-Delivery-System/Client"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

// Controller struct that holds a reference to the repository
type OrderController struct {
	Repo      *database.Repository
	ResClient *RestaurantClient.RestaurantClient
}

// NewController initializes the controller with the repository dependency
func NewController(repo *database.Repository, ResClient *RestaurantClient.RestaurantClient) *OrderController {
	return &OrderController{
		Repo:      repo,
		ResClient: ResClient,
	}
}

func (orderCtrl *OrderController) UpdateOrderStatus(c *gin.Context) {
	var OrderStatus payload.Order

	if err := c.ShouldBindJSON(&OrderStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order model.Order

	err := orderCtrl.Repo.GetOrder(&order, OrderStatus.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, "Order not found")
		return
	}

	if OrderStatus.DeliveryDriverID != 0 {
		order.DeliveryDriverID = OrderStatus.DeliveryDriverID

	}
	log.Printf("time: %v", order.Time)
	order.Time = time.Now()
	log.Printf("uptime: %v", order.Time)
	order.OrderStatus = OrderStatus.OrderStatus
	if err := orderCtrl.Repo.Update(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (orderCtrl *OrderController) GetOrders(c *gin.Context, UserType string) {

	var OrderNFilter payload.CombineOrderFilter
	if err := c.ShouldBindJSON(&OrderNFilter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order []model.Order
	var err error

	if UserType == "user" {
		err = orderCtrl.Repo.GetOrders(&order, int(OrderNFilter.UserId), OrderNFilter.Filter.ColumnName, OrderNFilter.Filter.OrderDirection, "user_id")
	} else if UserType == "restaurant" {
		err = orderCtrl.Repo.GetOrders(&order, int(OrderNFilter.RestaurantId), OrderNFilter.Filter.ColumnName, OrderNFilter.Filter.OrderDirection, "restaurant_id")

	} else if UserType == "delivery driver" {
		err = orderCtrl.Repo.GetOrders(&order, int(OrderNFilter.RestaurantId), OrderNFilter.Filter.ColumnName, OrderNFilter.Filter.OrderDirection, "delivery_driver")
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
	log.Printf("servie %s", ServiceType)
	if !exists {
		c.JSON(http.StatusBadRequest, "userId id does not exist")
		return
	}
	if ServiceType != "User" {
		c.JSON(http.StatusBadRequest, "Only user can place order")
		return
	}

	var CombineOrderItem payload.CombineOrderItem
	if err := c.ShouldBindJSON(&CombineOrderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var GetItem payload.GetItems
	GetItem.RestaurantId = CombineOrderItem.RestaurantId
	GetItem.ColumnName = "restaurant_id"
	GetItem.OrderType = "asc"

	var items []payload.Items
	items, err := orderCtrl.ResClient.GetItems(GetItem)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Error while getting items from the restaurant", "Error", err.Error())
		return
	}

	if len(items) == 0 {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "No items found in the restaurant", "", nil)
		return
	}

	totalBill, err := orderCtrl.calculateBill(CombineOrderItem, items)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	order := orderCtrl.createOrderObj(CombineOrderItem, totalBill)
	err = orderCtrl.Repo.PlaceOrder(&order, &CombineOrderItem)

	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Error while creating order", "Error", err.Error())
		return
	}
	utils.GenerateResponse(http.StatusOK, c, "Message", fmt.Sprintf("Order created successfully with total bill: %v", totalBill), "", nil)
}

func (orderCtrl *OrderController) ViewOrderDetails(c *gin.Context) {
	var orderId payload.ID

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
	var orderId payload.ID
	if err := c.ShouldBindJSON(&orderId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order model.Order
	if err := orderCtrl.Repo.GetOrder(&order, orderId.OrderID); err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Message", "Order not found", "Error", err.Error())
		return
	}

	var GetItem payload.GetItems
	GetItem.RestaurantId = order.RestaurantID
	GetItem.ColumnName = "restaurant_id"
	GetItem.OrderType = "asc"

	var items []payload.Items
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

	invoice := orderCtrl.createInvoice(order, orderItems, items)

	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}

func (orderCtrl *OrderController) createInvoice(order model.Order, orderItems []model.OrderItem, items []payload.Items) gin.H {
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
func (orderCtrl *OrderController) createOrderObj(order payload.CombineOrderItem, bill float64) model.Order {
	return model.Order{
		OrderStatus:  "order placed",
		UserId:       order.UserId,
		RestaurantID: order.RestaurantId,
		TotalBill:    bill,
	}
}
func (orderCtrl *OrderController) calculateBill(CombineOrderItem payload.CombineOrderItem, items []payload.Items) (float64, error) {
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
			return 0.0, fmt.Errorf("item with ID %d not found", orderedItem.ItemId)
		}

		totalBill += ItemPrice * float64(orderedItem.Quantity)
	}

	return totalBill, nil
}

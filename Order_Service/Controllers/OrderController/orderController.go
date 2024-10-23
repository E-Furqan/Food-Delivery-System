package OrderControllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

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
	log.Print("13111")
	var OrderStatus payload.Order

	if err := c.ShouldBindJSON(&OrderStatus); err != nil {
		log.Print("13")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order model.Order

	err := orderCtrl.Repo.GetOrder(&order, int(OrderStatus.OrderID))
	if err != nil {
		log.Print("23")
		c.JSON(http.StatusNotFound, "Order not found")
		return
	}

	if err := orderCtrl.Repo.Update(&order, OrderStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Print("33")
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
		err = orderCtrl.Repo.GetOrders(&order, int(OrderNFilter.UserId), OrderNFilter.Filter.ColumnName, OrderNFilter.Filter.OrderDirection, "user_id")
	} else {
		err = orderCtrl.Repo.GetOrders(&order, int(OrderNFilter.RestaurantId), OrderNFilter.Filter.ColumnName, OrderNFilter.Filter.OrderDirection, "restaurant_id")
		log.Print(order)
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

	var GetItem payload.GetItems
	GetItem.RestaurantId = CombineOrderItem.RestaurantId
	GetItem.ColumnName = "restaurant_id"
	GetItem.OrderType = "asc"

	items, err := orderCtrl.Client.GetItems(GetItem)
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
	order.UserId = CombineOrderItem.UserId
	order.RestaurantID = CombineOrderItem.RestaurantId
	order.TotalBill = totalBill
	order.OrderStatus = "order placed"

	err = orderCtrl.Repo.PlaceOrder(&order, &CombineOrderItem)

	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Error while creating order", "Error", err.Error())
		return
	}
	utils.GenerateResponse(http.StatusOK, c, "Message", "Order created successfully", "", nil)

	processOrder := orderCtrl.createProcessOrder(order)

	err = orderCtrl.Client.ProcessOrder(processOrder, false)

	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Error sending request to restaurant service", "Error", err.Error())
		return
	}

	utils.GenerateResponse(http.StatusOK, c, "Message", "Order Accepted by the restaurant successfully", "", nil)

}

func (orderCtrl *OrderController) StartScheduledOrderTask() {
	ticker := time.NewTicker(20 * time.Second)
	go func() {
		for {
			<-ticker.C
			orderCtrl.HandleProcessOrderTask()
		}
	}()
}

func (orderCtrl *OrderController) HandleProcessOrderTask() {
	var orders []model.Order

	if err := orderCtrl.Repo.FetchAllOrder(&orders); err != nil {
		log.Printf("Message: Error while getting all orders from order table. Error: %v", err.Error())
		return
	}
	var wg sync.WaitGroup
	currentTime := time.Now()
	for _, order := range orders {
		if strings.ToLower(order.OrderStatus) != "completed" && currentTime.Sub(order.Time) >= 50*time.Second {
			processOrder := orderCtrl.createProcessOrder(order)
			if err := orderCtrl.ProcessOrderBasedOnStatus(processOrder, order.OrderStatus); err != nil {
				log.Printf("Message: Error sending request to service. Error: %v", err.Error())
				log.Printf("order id :%v", order.OrderID)
			}

		}
	}
	wg.Wait()
	log.Print("All orders processed")
}

func (orderCtrl *OrderController) createProcessOrder(order model.Order) payload.ProcessOrder {
	return payload.ProcessOrder{
		OrderStatus: order.OrderStatus,
		ID: payload.ID{
			RestaurantId: order.RestaurantID,
			OrderID:      order.OrderID,
		},
	}
}

func (orderCtrl *OrderController) ProcessOrderBasedOnStatus(processOrder payload.ProcessOrder, status string) error {
	if orderCtrl.isRestaurantStatus(status) {
		return orderCtrl.Client.ProcessOrder(processOrder, false)
	}
	if orderCtrl.isUserStatus(status) {
		return orderCtrl.Client.ProcessOrder(processOrder, true)
	}
	return fmt.Errorf("status not recognized: %s", status)
}

func (orderCtrl *OrderController) isRestaurantStatus(status string) bool {
	for _, restaurantStatus := range payload.RestaurantOrderStatuses {
		if strings.EqualFold(status, restaurantStatus) {
			return true
		}
	}
	return false
}

func (orderCtrl *OrderController) isUserStatus(status string) bool {
	for _, userStatus := range payload.UserOrderStatuses {
		if strings.EqualFold(status, userStatus) {
			log.Print("User")
			return true
		}
	}
	return false
}

package RestaurantController

import (
	"log"
	"net/http"
	"strings"

	ClientPackage "github.com/E-Furqan/Food-Delivery-System/Client"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type RestaurantController struct {
	Repo   *database.Repository
	Client *ClientPackage.Client
}

func NewController(repo *database.Repository, client *ClientPackage.Client) *RestaurantController {
	return &RestaurantController{
		Repo:   repo,
		Client: client,
	}
}

func (ctrl *RestaurantController) Register(c *gin.Context) {

	var registrationData model.Restaurant

	if err := c.ShouldBindJSON(&registrationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.Repo.CreateRestaurant(&registrationData)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": pqErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, registrationData)
}

func (ctrl *RestaurantController) Login(c *gin.Context) {

	var input payload.Credentials
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error while binding": err.Error()})
		return
	}

	var Restaurant model.Restaurant
	err := ctrl.Repo.GetRestaurant("restaurant_email", input.Email, &Restaurant)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if Restaurant.Password != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	var RestaurantClaim payload.RestaurantClaim
	RestaurantClaim.ClaimId = Restaurant.RestaurantId
	RestaurantClaim.ServiceType = "Restaurant"

	tokens, err := ctrl.Client.GenerateResponse(RestaurantClaim)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access token":  tokens.AccessToken,
		"refresh token": tokens.RefreshToken,
		"expires at":    tokens.Expiration,
	})
}

func (ctrl *RestaurantController) GetAllRestaurants(c *gin.Context) {

	var restaurants []model.Restaurant

	if err := ctrl.Repo.GetAllRestaurants(&restaurants); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, restaurants)
}

func (ctrl *RestaurantController) UpdateRestaurantStatus(c *gin.Context) {
	RestaurantID, err := utils.Verification(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant not authenticated"})
		return
	}

	var Restaurant model.Restaurant
	err = ctrl.Repo.GetRestaurant("restaurant_id", RestaurantID, &Restaurant)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var input payload.Input

	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error while binding": err})
		return
	}

	if err = ctrl.Repo.UpdateRestaurantStatus(&Restaurant, input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error while updating": err})
		return
	}

	c.JSON(http.StatusOK, "Restaurant status updated")
}

func (ctrl *RestaurantController) ViewMenu(c *gin.Context) {

	var Items []model.Item
	var combinedInput payload.CombinedInput

	if err := c.ShouldBindJSON(&combinedInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error binding": err.Error()})
		return
	}
	var Restaurant model.Restaurant
	err := ctrl.Repo.GetRestaurant("restaurant_id", combinedInput.RestaurantId, &Restaurant)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant does not exist"})
		return
	}

	Items, err = ctrl.Repo.LoadItems(combinedInput.RestaurantId, combinedInput.ColumnName, combinedInput.OrderType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error load item": err.Error()})
		return
	}

	if len(Items) <= 0 {
		c.JSON(http.StatusInternalServerError, "No items present in the restaurant")
		return
	}

	c.JSON(http.StatusOK, Items)
}

func (ctrl *RestaurantController) ProcessOrder(c *gin.Context) {

	RestaurantID, err := utils.Verification(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant not authenticated"})
		return
	}
	// var Restaurant model.Restaurant
	// err := ctrl.Repo.GetRestaurant("restaurant_id", RestaurantID, &Restaurant)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant does not exists"})
	// 	return
	// }

	var order payload.OrderDetails

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, "Error while binding order status")
		return
	}

	OrderDetails, err := ctrl.Client.ViewOrdersDetails(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if OrderDetails.RestaurantId != RestaurantID {
		log.Printf("order %s res %v", OrderDetails.OrderStatus, RestaurantID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order is not for your restaurant"})
		return
	}

	if strings.ToLower(order.OrderStatus) == "cancelled" {
		OrderDetails.OrderStatus = "Cancelled"
		if err := ctrl.Client.ProcessOrder(*OrderDetails); err != nil {
			utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Post request failed", "Error", err.Error())
			return
		}
		utils.GenerateResponse(http.StatusOK, c, "Message", "Post request successful", "", "")
		return
	}

	orderTransitions := payload.GetOrderTransitions()

	if strings.ToLower(OrderDetails.OrderStatus) == "order placed" {
		var restaurant model.Restaurant
		err := ctrl.Repo.GetRestaurant("restaurant_id", OrderDetails.RestaurantId, &restaurant)
		if err != nil {
			OrderDetails.OrderStatus = "Cancelled"
			c.JSON(http.StatusNotFound, "Restaurant not found")
			c.JSON(http.StatusNotFound, OrderDetails)
			log.Printf("restaurant not found cancel")
			return
		}

		if strings.ToLower(restaurant.RestaurantStatus) == "closed" {
			OrderDetails.OrderStatus = "Cancelled"
			c.JSON(http.StatusBadRequest, "Restaurant is closed")
			c.JSON(http.StatusBadRequest, OrderDetails)
			log.Printf("restaurant is close")
			return
		}
	}

	if newStatus, exists := orderTransitions[OrderDetails.OrderStatus]; exists {
		OrderDetails.OrderStatus = newStatus
	}

	if err := ctrl.Client.ProcessOrder(*OrderDetails); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Post request failed", "Error", err.Error())
		return
	}

	utils.GenerateResponse(http.StatusOK, c, "Message", "Post request successful", "", "")
}

// // remove
// func (ctrl *RestaurantController) CancelOrder(c *gin.Context) {
// 	RestaurantID, exists := c.Get("RestaurantID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant not authenticated"})
// 		return
// 	}
// 	RestaurantID, ok := RestaurantID.(uint)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Email address"})
// 		return
// 	}

// 	var Restaurant model.Restaurant
// 	err := ctrl.Repo.GetRestaurant("restaurant_id", RestaurantID, &Restaurant)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
// 		return
// 	}

// 	var input payload.OrderDetails

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Binding input data failed", "Error", err.Error())
// 		return
// 	}

// 	if input.RestaurantId != Restaurant.RestaurantId {

// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"Message": "You are not authorized to cancel this order as it belongs to a different restaurant",
// 		})
// 		return
// 	}

// 	input.OrderStatus = "Cancelled"
// 	utils.GenerateResponse(http.StatusOK, c, "Message", "Order cancelled successfully", "Order details", input)

// 	if err := ctrl.Client.ProcessOrder(input); err != nil {
// 		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Post request failed", "Error", err.Error())
// 		return
// 	}

// 	utils.GenerateResponse(http.StatusOK, c, "Message", "Post request successful", "", nil)
// }

func (ctrl *RestaurantController) RefreshToken(c *gin.Context) {

	var input payload.RefreshToken

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var refreshClaim payload.RefreshToken
	refreshClaim.RefreshToken = input.RefreshToken
	refreshClaim.ServiceType = "Restaurant"
	tokens, err := ctrl.Client.RefreshToken(refreshClaim)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access token":  tokens.AccessToken,
		"refresh token": tokens.RefreshToken,
		"expires at":    tokens.Expiration,
	})
}

func (ctrl *RestaurantController) ViewRestaurantOrders(c *gin.Context) {
	RestaurantID, err := utils.Verification(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant not authenticated"})
		return
	}

	var Restaurant model.Restaurant
	err = ctrl.Repo.GetRestaurant("restaurant_id", RestaurantID, &Restaurant)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant does not exists"})
		return
	}
	var restaurantId payload.Input

	restaurantId.RestaurantId = Restaurant.RestaurantId
	Orders, err := ctrl.Client.ViewRestaurantOrders(restaurantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error order": err.Error()})
		return
	}
	var filter payload.OrderFilter
	if err := c.ShouldBindJSON(filter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error order": err.Error()})
		return
	}

	var filteredOrders []payload.OrderDetails

	// Filter the orders based on OrderStatus
	for _, order := range *Orders {
		if filter.Filter == "all" || filter.Filter == "" || order.OrderStatus == filter.Filter {
			filteredOrders = append(filteredOrders, order)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Restaurant orders: ": filteredOrders,
	})
}

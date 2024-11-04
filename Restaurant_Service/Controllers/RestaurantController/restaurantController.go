package RestaurantController

import (
	"log"
	"net/http"
	"strings"

	"github.com/E-Furqan/Food-Delivery-System/Client/AuthClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type RestaurantController struct {
	Repo        *database.Repository
	OrderClient *OrderClient.OrderClient
	AuthClient  *AuthClient.AuthClient
}

func NewController(repo *database.Repository, OrderClient *OrderClient.OrderClient, AuthClient *AuthClient.AuthClient) *RestaurantController {
	return &RestaurantController{
		Repo:        repo,
		OrderClient: OrderClient,
		AuthClient:  AuthClient,
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

	var input model.Credentials
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

	var RestaurantClaim model.RestaurantClaim
	RestaurantClaim.ClaimId = Restaurant.RestaurantId
	RestaurantClaim.Role = "Restaurant"

	tokens, err := ctrl.AuthClient.GenerateResponse(RestaurantClaim)
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

	var input model.Input

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
	var combinedInput model.CombinedInput

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

	var order model.OrderDetails

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, "Error while binding order status")
		return
	}

	OrderDetails, err := ctrl.OrderClient.ViewOrdersDetails(order)
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
		if err := ctrl.OrderClient.ProcessOrder(*OrderDetails); err != nil {
			utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Post request failed", "Error", err.Error())
			return
		}
		utils.GenerateResponse(http.StatusOK, c, "Message", "Post request successful", "", "")
		return
	}

	orderTransitions := model.GetOrderTransitions()

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

	if err := ctrl.OrderClient.ProcessOrder(*OrderDetails); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Post request failed", "Error", err.Error())
		return
	}

	utils.GenerateResponse(http.StatusOK, c, "Message", "Order status updated", "", "")
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
	var restaurantId model.Input

	restaurantId.RestaurantId = Restaurant.RestaurantId
	Orders, err := ctrl.OrderClient.ViewRestaurantOrders(restaurantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error order": err.Error()})
		return
	}
	var filter model.OrderFilter
	if err := c.ShouldBindJSON(filter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error order": err.Error()})
		return
	}

	var filteredOrders []model.OrderDetails

	// Filter the orders based on OrderStatus
	for _, order := range *Orders {
		if filter.Filter == "all" || filter.Filter == "" || strings.EqualFold(order.OrderStatus, filter.Filter) {
			filteredOrders = append(filteredOrders, order)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Restaurant orders: ": filteredOrders,
	})
}

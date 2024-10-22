package RestaurantController

import (
	"log"
	"net/http"
	"time"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	access_token, refresh_token, err := utils.GenerateTokens(Restaurant.RestaurantEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access token":  access_token,
		"refresh token": refresh_token,
		"expires at":    time.Now().Add(24 * time.Hour).Unix(),
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

	email, exists := c.Get("Email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	email, ok := email.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Email address"})
		return
	}

	var Restaurant model.Restaurant
	err := ctrl.Repo.GetRestaurant("restaurant_email", email, &Restaurant)
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
	log.Print("combined Input :", combinedInput)
	Items, err := ctrl.Repo.LoadItems(combinedInput.RestaurantId, combinedInput.ColumnName, combinedInput.OrderType)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error load item": err.Error()})
		return
	}

	log.Print("Fetched items:", Items)
	c.JSON(http.StatusOK, Items)
}

func (ctrl *RestaurantController) ProcessOrder(c *gin.Context) {
	var order payload.ProcessOrder

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, "Error while binding order status")
		return
	}

	orderTransitions := payload.GetOrderTransitions()

	if order.OrderStatus == "order placed" {
		var restaurant model.Restaurant
		err := ctrl.Repo.GetRestaurant("restaurant_id", order.RestaurantId, &restaurant)
		if err != nil {
			order.OrderStatus = "Cancelled"
			c.JSON(http.StatusNotFound, "Restaurant not found")
			c.JSON(http.StatusNotFound, order)
			return
		}

		if restaurant.RestaurantStatus == "closed" || restaurant.RestaurantStatus == "Closed" {
			order.OrderStatus = "Cancelled"
			c.JSON(http.StatusBadRequest, "Restaurant is closed")
			c.JSON(http.StatusBadRequest, order)
			return
		}
	}

	if newStatus, exists := orderTransitions[order.OrderStatus]; exists {
		order.OrderStatus = newStatus
	}

	if err := ctrl.Client.ProcessOrder(order); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Post request failed", "Error", err.Error())
		return
	}

	utils.GenerateResponse(http.StatusOK, c, "Message", "Post request successful", "", "")

}

func (ctrl *RestaurantController) CancelOrder(c *gin.Context) {
	email, exists := c.Get("Email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant not authenticated"})
		return
	}
	email, ok := email.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Email address"})
		return
	}

	var Restaurant model.Restaurant
	err := ctrl.Repo.GetRestaurant("restaurant_email", email, &Restaurant)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var input payload.ProcessOrder

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Binding input data failed", "Error", err.Error())
		return
	}

	if input.RestaurantId != Restaurant.RestaurantId {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "You are not authorized to cancel this order as it belongs to a different restaurant",
		})
		return
	}

	input.OrderStatus = "Cancelled"
	utils.GenerateResponse(http.StatusOK, c, "Message", "Order cancelled successfully", "Order details", input)

	if err := ctrl.Client.ProcessOrder(input); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Post request failed", "Error", err.Error())
		return
	}

	utils.GenerateResponse(http.StatusOK, c, "Message", "Post request successful", "", nil)
}

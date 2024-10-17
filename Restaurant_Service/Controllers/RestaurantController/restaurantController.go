package controllers

import (
	"net/http"
	"time"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// Controller struct that holds a reference to the repository
type RestaurantController struct {
	Repo *database.Repository
}

// NewController initializes the controller with the repository dependency
func NewController(repo *database.Repository) *RestaurantController {
	return &RestaurantController{Repo: repo}
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

	// Respond with the created user data
	c.JSON(http.StatusCreated, registrationData)
}

func (ctrl *RestaurantController) Login(c *gin.Context) {

	var input payload.Credentials
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var Restaurant model.Restaurant
	err := ctrl.Repo.GetRestaurant("RestaurantEmail", input.Email, &Restaurant)
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

func (ctrl *RestaurantController) ViewMenu(c *gin.Context) {

	var Items []model.Item
	var OrderInfo payload.Order
	var input payload.Input

	if err := c.ShouldBindJSON(&OrderInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Items, err := ctrl.Repo.LoadItemsInOrder(input.RestaurantId, OrderInfo.ColumnName, OrderInfo.OrderType)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Items)
}

func (ctrl *RestaurantController) AddItemOfRestaurantMenu(c *gin.Context) {
	email, exists := c.Get("Email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	email, ok := email.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Email address"})
		return
	}

	var Restaurant model.Restaurant
	err := ctrl.Repo.GetRestaurant("RestaurantEmail", email, &Restaurant)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var NewItemData model.Item

	if err = c.ShouldBindJSON(&NewItemData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err = ctrl.Repo.AddItemToRestaurantMenu(Restaurant.RestaurantId, NewItemData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, "Item added to the restaurant menu")
}

func (ctrl *RestaurantController) DeleteItemsOfRestaurantMenu(c *gin.Context) {
	email, exists := c.Get("Email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	email, ok := email.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Email address"})
		return
	}

	var Restaurant model.Restaurant
	err := ctrl.Repo.GetRestaurant("RestaurantEmail", email, &Restaurant)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var DeleteItemId payload.Input

	if err = c.ShouldBindJSON(&DeleteItemId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err = ctrl.Repo.RemoveItemFromRestaurantMenu(Restaurant.RestaurantId, DeleteItemId.ItemId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, "Item deleted from the restaurant menu")
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
	email, ok := email.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Email address"})
		return
	}

	var Restaurant model.Restaurant
	err := ctrl.Repo.GetRestaurant("RestaurantEmail", email, &Restaurant)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var input payload.Input

	if err = c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err = ctrl.Repo.UpdateRestaurantStatus(&Restaurant, input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, Restaurant)

}

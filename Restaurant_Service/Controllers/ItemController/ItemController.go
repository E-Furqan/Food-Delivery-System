package ItemController

import (
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
)

type ItemController struct {
	Repo *database.Repository
}

func NewController(repo *database.Repository) *ItemController {
	return &ItemController{Repo: repo}
}

func (ItemController *ItemController) ViewMenu(c *gin.Context) {

	var Items []model.Item
	var combinedInput payload.CombinedInput

	if err := c.ShouldBindJSON(&combinedInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error binding": err.Error()})
		return
	}
	Items, err := ItemController.Repo.LoadItemsInOrder(combinedInput.RestaurantId, combinedInput.ColumnName, combinedInput.OrderType)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error load item": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Items)
}

func (ItemController *ItemController) AddItemItRestaurantMenu(c *gin.Context) {
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
	err := ItemController.Repo.GetRestaurant("restaurant_email", email, &Restaurant)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var NewItemData model.Item

	if err = c.ShouldBindJSON(&NewItemData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err = ItemController.Repo.AddItemToRestaurantMenu(Restaurant.RestaurantId, NewItemData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Restaurant)
}

func (ItemController *ItemController) DeleteItemsOfRestaurantMenu(c *gin.Context) {
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
	err := ItemController.Repo.GetRestaurant("restaurant_email", email, &Restaurant)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var DeleteItemId payload.Input

	if err = c.ShouldBindJSON(&DeleteItemId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err = ItemController.Repo.RemoveItemFromRestaurantMenu(Restaurant.RestaurantId, DeleteItemId.ItemId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, "Item deleted from the restaurant menu")
}

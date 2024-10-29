package ItemController

import (
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

type ItemController struct {
	Repo *database.Repository
}

func NewController(repo *database.Repository) *ItemController {
	return &ItemController{Repo: repo}
}

func (ItemController *ItemController) AddItemsInMenu(c *gin.Context) {

	RestaurantID, err := utils.Verification(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant not authenticated"})
		return
	}
	var Restaurant model.Restaurant
	err = ItemController.Repo.GetRestaurant("restaurant_id", RestaurantID, &Restaurant)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var NewItemData model.Item

	if err = c.ShouldBindJSON(&NewItemData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	NewItemData.RestaurantId = Restaurant.RestaurantId
	if err = ItemController.Repo.AddItemToRestaurantMenu(NewItemData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Item added to menu successfully")
}

func (ItemController *ItemController) DeleteItemsFromMenu(c *gin.Context) {

	RestaurantID, err := utils.Verification(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant not authenticated"})
		return
	}

	var Restaurant model.Restaurant
	err = ItemController.Repo.GetRestaurant("restaurant_id", RestaurantID, &Restaurant)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var DeleteItemId model.Input

	if err = c.ShouldBindJSON(&DeleteItemId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err = ItemController.Repo.RemoveItem(Restaurant.RestaurantId, DeleteItemId.ItemId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, "Item deleted from the restaurant menu")
}

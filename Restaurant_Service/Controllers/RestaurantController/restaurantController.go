package controllers

import (
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
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

package UserControllers

import (
	"fmt"
	"net/http"
	"time"

	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	model "github.com/E-Furqan/Food-Delivery-System/models"
	"github.com/E-Furqan/Food-Delivery-System/payload"
	"github.com/E-Furqan/Food-Delivery-System/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// Controller struct that holds a reference to the repository
type Controller struct {
	Repo *database.Repository
}

// NewController initializes the controller with the repository dependency
func NewController(repo *database.Repository) *Controller {
	return &Controller{Repo: repo}
}

func (ctrl *Controller) Register(c *gin.Context) {

	var registrationData model.User

	if err := c.ShouldBindJSON(&registrationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.Repo.CreateUser(&registrationData)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			// Return only the Message field from the error
			c.JSON(http.StatusInternalServerError, gin.H{"error": pqErr.Message})
		} else {
			// Return a generic error message if it's not a pq.Error
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, registrationData)
}

func (ctrl *Controller) Login(c *gin.Context) {

	var input payload.Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	err := ctrl.Repo.FindUser(input.Username, &user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.Password != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	access_token, refresh_token, err := utils.GenerateTokens(user.Username)
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

func (ctrl *Controller) GetUser(c *gin.Context) {

	var userData []model.User
	var OrderInfo payload.Order

	if err := c.ShouldBindJSON(&OrderInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, err := ctrl.Repo.PreloadInOrder(OrderInfo.ColumnName, OrderInfo.OrderType)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userData)

}

func (ctrl *Controller) UpdateUser(c *gin.Context) {
	var user model.User

	// Retrieve the username from the context
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	// Perform type assertion to string
	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username type"})
		return
	}

	// Fetch the user by username
	err := ctrl.Repo.FindUser(usernameStr, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User not found %v %s", err, usernameStr)})
		return
	}

	var update_user model.User
	err1 := c.ShouldBindJSON(&update_user)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}

	err = ctrl.Repo.Update(&user, &update_user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (ctrl *Controller) DeleteUser(c *gin.Context) {
	// Retrieve the username from the context
	usernameValue, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	username, ok := usernameValue.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username type"})
		return
	}

	user := model.User{}
	err := ctrl.Repo.FindUser(username, &user)

	// Fetch the user by username
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete the user
	if err := ctrl.Repo.DeleteUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})

}

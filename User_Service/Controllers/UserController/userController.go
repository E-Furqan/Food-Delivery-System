package UserControllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Repo *database.Repository
}

func NewController(repo *database.Repository) *Controller {
	return &Controller{Repo: repo}
}

func (ctrl *Controller) Register(c *gin.Context) {

	var registrationData model.User

	if err := c.ShouldBindJSON(&registrationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(registrationData.Roles) > 0 && registrationData.ActiveRole == "" {
		var role model.Role
		if err := ctrl.Repo.GetRole(registrationData.Roles[0].RoleId, &role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Role not found"})
			return
		}
		registrationData.ActiveRole = role.RoleType
		log.Print("active role set")
	}

	err := ctrl.Repo.CreateUser(&registrationData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, registrationData)
}

func (ctrl *Controller) Login(c *gin.Context) {

	var input payload.Credentials
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	err := ctrl.Repo.GetUser("username", input.Username, &user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.Password != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	access_token, refresh_token, err := utils.GenerateTokens(user.Username, user.ActiveRole)
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

func (ctrl *Controller) GetUsers(c *gin.Context) {

	activeRole, exists := c.Get("activeRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if activeRole != "Admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You do not have the privileges to view users."})
		return
	}

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

	usernameValue, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "User not authenticated"})
		return
	}
	username, ok := usernameValue.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username type"})
		return
	}

	user := model.User{}
	err := ctrl.Repo.GetUser("username", username, &user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateUserData model.User
	err = c.ShouldBindJSON(&updateUserData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.Repo.DeleteUserRoleInfo(user.UserId, "user_user_id"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var roleType string
	for _, role := range updateUserData.Roles {
		if role.RoleId == 3 || role.RoleType == "Admin" {
			roleType = "Admin"
			break
		} else {
			roleType = role.RoleType
		}
	}

	user.Roles = updateUserData.Roles
	user.ActiveRole = roleType
	err = ctrl.Repo.Update(&user, &updateUserData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}

	for _, role := range updateUserData.Roles {
		if err := ctrl.Repo.AddUserRole(user.UserId, role.RoleId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add roles to user"})
			return
		}
	}
	log.Print(updateUserData.Roles)
	log.Print(user.Roles)

	c.JSON(http.StatusCreated, user)
}

func (ctrl *Controller) DeleteUser(c *gin.Context) {
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
	err := ctrl.Repo.GetUser("username", username, &user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := ctrl.Repo.DeleteUserRoleInfo(user.UserId, "user_user_id"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.Repo.DeleteUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (ctrl *Controller) Profile(c *gin.Context) {

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username type"})
		return
	}

	var user model.User

	err := ctrl.Repo.GetUser("username", usernameStr, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User not found %v %s", err, usernameStr)})
		return
	}

	c.JSON(http.StatusFound, user)
}

func (ctrl *Controller) SearchForUser(c *gin.Context) {
	role, exists := c.Get("activeRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	if role != "Admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You do not have the privileges to Search for users."})
		return
	}

	var input payload.UserSearch
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	err = ctrl.Repo.GetUser(input.ColumnName, input.SearchParameter, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Error: %v \nUser not found: %s", err, input.SearchParameter)})
		return
	}

	c.JSON(http.StatusFound, user)
}

func (ctrl *Controller) SwitchRole(c *gin.Context) {

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username type"})
		return
	}

	var user model.User
	err := ctrl.Repo.GetUser("username", usernameStr, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User not found %v %s", err, usernameStr)})
		return
	}

	var RoleSwitch payload.RoleSwitch
	err = c.ShouldBindJSON(&RoleSwitch)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var roleExists bool
	var newRole model.Role
	for _, role := range user.Roles {
		if role.RoleId == RoleSwitch.NewRoleID {
			roleExists = true
			newRole = role
			break
		}
	}

	if !roleExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found in user's roles"})
		return
	}

	user.ActiveRole = newRole.RoleType

	if err := ctrl.Repo.UpdateUserActiveRole(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user active role"})
		return
	}

	access_token, refresh_token, err := utils.GenerateTokens(user.Username, user.ActiveRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":              user,
		"new access token":  access_token,
		"new refresh token": refresh_token,
		"expires at":        time.Now().Add(24 * time.Hour).Unix(),
	})

}

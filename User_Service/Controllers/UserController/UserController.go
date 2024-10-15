package UserControllers

import (
	"fmt"
	"log"
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

	// Bind JSON data to the User struct
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
	var roleIds []uint // Create a slice to hold the role IDs

	// Extract Role IDs from user.Roles
	for _, role := range user.Roles {
		roleIds = append(roleIds, role.RoleId) // Assuming role.RoleId is of type uint
	}

	access_token, refresh_token, err := utils.GenerateTokens(user.Username, roleIds)
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

	// Retrieve the slice of Role IDs from the context
	roleIdsValue, exists := c.Get("roleId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	roleIds, ok := roleIdsValue.([]uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role id type"})
		return
	}

	// Check if any of the role IDs is for Admin
	var isAdmin bool
	var role model.Role

	for _, roleId := range roleIds {
		err := ctrl.Repo.GetRole(roleId, &role)
		if err == nil && role.RoleType == "Admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You do not have the privileges to view roles."})
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
	// Retrieve the username from the context
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

	// Fetch the user by username
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

	// First, delete the associations in the user_roles table
	if err := ctrl.Repo.DeleteUserRoleInfo(user.UserId, "user_user_id"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Fetch the user by username
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var roleType string
	for _, role := range updateUserData.Roles {
		if role.RoleId == 3 || role.RoleType == "Admin" {
			roleType = "Admin"
			break
		} else {
			roleType = role.RoleType // Save the first other role type encountered
		}
	}

	user.Roles = updateUserData.Roles
	user.ActiveRole = roleType
	err = ctrl.Repo.Update(&user, &updateUserData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}

	// Recreate associations with new roles
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
	err := ctrl.Repo.GetUser("username", username, &user)

	// Fetch the user by username
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// First, delete the associations in the user_roles table
	if err := ctrl.Repo.DeleteUserRoleInfo(user.UserId, "user_user_id"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

func (ctrl *Controller) Profile(c *gin.Context) {

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

	var user model.User
	// Fetch the user by username
	err := ctrl.Repo.GetUser("username", usernameStr, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User not found %v %s", err, usernameStr)})
		return
	}

	c.JSON(http.StatusFound, user)
}

// admin only for admin
func (ctrl *Controller) SearchForUser(c *gin.Context) {
	var input payload.UserSearch
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	// Fetch the user by username
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

	// Check if the new role exists in the user's roles
	var roleExists bool
	var newRole model.Role
	for _, role := range user.Roles {
		if role.RoleId == RoleSwitch.NewRoleID { // Check if role exists in user's roles
			roleExists = true
			newRole = role
			break
		}
	}

	// If the role doesn't exist, return an error
	if !roleExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found in user's roles"})
		return
	}

	user.ActiveRole = newRole.RoleType
	// Save the updated user to the database
	if err := ctrl.Repo.UpdateUserActiveRole(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user active role"})
		return
	}

	c.JSON(http.StatusOK, user)

}

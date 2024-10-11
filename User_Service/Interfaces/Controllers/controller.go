package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	database "github.com/E-Furqan/Food-Delivery-System/Interfaces/Repositories"
	environmentvariable "github.com/E-Furqan/Food-Delivery-System/enviorment_variable"
	model "github.com/E-Furqan/Food-Delivery-System/models"
	"github.com/E-Furqan/Food-Delivery-System/utils"
	"github.com/gin-gonic/gin"
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

	var reg_data model.User

	if err := c.ShouldBindJSON(&reg_data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := ctrl.Repo.CheckUserExistence(reg_data.Username, reg_data.Email, int(reg_data.Phone_number))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user existence"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, phone number, or email already exists"})
		return
	}

	var role model.Role
	role_made := ctrl.Check_If_Role_Exist(reg_data.Role_id, c, &role)
	if !role_made {
		return
	}

	err = ctrl.Repo.CreateUser(&reg_data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	log.Printf("User created successfully: %+v", reg_data)

	reg_data, err = ctrl.Repo.LoadUserWithRole(reg_data.User_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": reg_data.User_id})
		return
	}

	c.JSON(http.StatusCreated, reg_data)
}

func (ctrl *Controller) Check_If_Role_Exist(Role_id string, c *gin.Context, role *model.Role) bool {
	// Check if the role exists

	err := ctrl.Repo.Find_Role_By_Role_Id(Role_id, role)
	if err != nil {

		if Role_id == "1" {

			role.Role_id = Role_id
			role.Role_type = "Customer"

			ctrl.Repo.CreateRole(role)

		} else if Role_id == "2" {

			role.Role_id = Role_id
			role.Role_type = "Delivery driver"

			ctrl.Repo.CreateRole(role)

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
			return false
		}

	}
	return true
}

func (ctrl *Controller) Login(c *gin.Context) {

	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	err := ctrl.Repo.Find_User_By_Username(input.Username, &user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.Password != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_at": time.Now().Add(24 * time.Hour).Unix(),
	})

}

func (ctrl *Controller) Get_user(c *gin.Context) {

	var user_data []model.User
	user_data, err := ctrl.Repo.Preload_in_order()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user_data)

}
func (ctrl *Controller) Get_role(c *gin.Context) {
	var user_data []model.Role
	user_data, err := ctrl.Repo.Role_in_Asc_order()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user_data)
}

func (ctrl *Controller) Update_Role(c *gin.Context) {

	// Retrieve the username from the context
	usernameValue, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// // Define the input structure for binding
	// var input struct {
	// 	Role_id string `json:"role_id" binding:"required"`
	// }

	var update_data model.User
	var user model.User
	// Bind incoming JSON to input struct
	if err := c.ShouldBindJSON(&update_data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the Role_id exists in the database
	var role model.Role
	role_made := ctrl.Check_If_Role_Exist(update_data.Role_id, c, &role)
	if !role_made {
		return
	}

	// Ensure that the username is a string
	username, ok := usernameValue.(string) // Type assertion for usernameValue
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username type"})
		return
	}

	// Fetch the user by username
	err := ctrl.Repo.Find_User_By_Username(username, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update the user's role_id
	if err := ctrl.Repo.Update(&user, &update_data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully", "role_id": user.Role_id})
}

func (ctrl *Controller) Update_user(c *gin.Context) {
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
	err := ctrl.Repo.Find_User_By_Username(usernameStr, &user)
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

	err = ctrl.Repo.Preload_Role_first(&user, int(user.User_id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user with role"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (ctrl *Controller) Delete_user(c *gin.Context) {
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
	err := ctrl.Repo.Find_User_By_Username(username, &user)

	// Fetch the user by username
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete the user
	if err := ctrl.Repo.Delete_user(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})

}

func (ctrl *Controller) Delete_role(c *gin.Context) {

	Admin := environmentvariable.Get_env("ADMIN")
	Admin_password := environmentvariable.Get_env("ADMIN_PASS")

	// Define the input structure for binding
	var input struct {
		Username string `json:"username_admin" binding:"required"`
		Password string `json:"password_admin" binding:"required"`
		Role_id  string `json:"role_id" binding:"required"`
	}

	// Bind incoming JSON to input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var role model.Role
	if input.Username == Admin && input.Password == Admin_password {
		// Fetch the role by role id
		err := ctrl.Repo.Find_Role_By_Role_Id(input.Role_id, &role)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		// Delete the role
		if err := ctrl.Repo.Delete_role(&role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the role"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
	}

}

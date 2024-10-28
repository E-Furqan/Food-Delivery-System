package UserControllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	ClientPackage "github.com/E-Furqan/Food-Delivery-System/Client"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Repo   *database.Repository
	Client *ClientPackage.Client
}

func NewController(repo *database.Repository, client *ClientPackage.Client) *Controller {
	return &Controller{
		Repo:   repo,
		Client: client,
	}
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
	var UserClaim payload.UserClaim
	UserClaim.Username = user.Username
	UserClaim.ActiveRole = user.ActiveRole
	UserClaim.ServiceType = "User"
	token, err := ctrl.Client.GenerateResponse(UserClaim)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access token":  token.AccessToken,
		"refresh token": token.RefreshToken,
		"expires at":    token.Expiration,
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

	username, err := utils.VerificationUsername(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user := model.User{}
	err = ctrl.Repo.GetUser("username", username, &user)

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
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
	username, err := utils.VerificationUsername(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user := model.User{}
	err = ctrl.Repo.GetUser("username", username, &user)

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

	var UserClaim payload.UserClaim
	UserClaim.Username = user.Username
	UserClaim.ActiveRole = user.ActiveRole
	UserClaim.ServiceType = "User"
	token, err := ctrl.Client.GenerateResponse(UserClaim)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access token":  token.AccessToken,
		"refresh token": token.RefreshToken,
		"expires at":    token.Expiration,
	})
}

// func (ctrl *Controller) ProcessOrder(c *gin.Context) {
// 	var order payload.ProcessOrder

// 	if err := c.ShouldBindJSON(&order); err != nil {
// 		c.JSON(http.StatusBadRequest, "Error while binding order status")
// 		return
// 	}

// 	if order.OrderStatus == "Cancelled" {
// 		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Order Cancelled", "order", order)
// 		return
// 	}

// 	orderTransitions := payload.GetOrderTransitions()
// 	if order.OrderStatus == "Waiting For Delivery Driver" {
// 		var driver model.User
// 		err := ctrl.Repo.GetDeliveryDriver(&driver)
// 		if err != nil {
// 			utils.GenerateResponse(http.StatusNotFound, c, "Message", "Delivery driver not found", "Error", err.Error())
// 			return
// 		}

// 		if newStatus, exists := orderTransitions[order.OrderStatus]; exists {
// 			order.OrderStatus = newStatus
// 		}
// 		order.DeliverDriverID = driver.UserId

// 		if err := ctrl.Client.ProcessOrder(order); err != nil {
// 			utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Patch request failed", "Error", err.Error())
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{
// 			"Order":          order,
// 			"Deliver Driver": driver.UserId,
// 		})
// 		return

// 	} else if order.OrderStatus == "Delivered" {
// 		var driver model.User
// 		log.Print(order.DeliverDriverID)
// 		err := ctrl.Repo.GetUser("user_id", order.DeliverDriverID, &driver)
// 		if err != nil {
// 			utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Delivery driver not found", "Error", err.Error())
// 			return
// 		}

// 		driver.RoleStatus = "available"
// 		ctrl.Repo.UpdateRoleStatus(&driver)
// 	}

// 	if newStatus, exists := orderTransitions[order.OrderStatus]; exists {
// 		log.Print(newStatus)
// 		order.OrderStatus = newStatus
// 	}
// 	if err := ctrl.Client.ProcessOrder(order); err != nil {
// 		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Patch request failed", "Error", err.Error())
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"Order": order,
// 	})
// }

func (ctrl *Controller) ViewUserOrders(c *gin.Context) {
	username, err := utils.VerificationUsername(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant not authenticated"})
		return
	}

	var User model.User
	err = ctrl.Repo.GetUser("username", username, &User)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant does not exists"})
		return
	}
	var userId payload.ProcessOrder

	userId.UserID = User.UserId
	Orders, err := ctrl.Client.ViewUserOrders(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error order": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"User orders: ": Orders,
	})
}

func (ctrl *Controller) ProcessOrderUser(c *gin.Context) {

	username, err := utils.VerificationUsername(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user := model.User{}
	err = ctrl.Repo.GetUser("username", username, &user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var order payload.ProcessOrder

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, "Error while binding order status")
		return
	}
	OrderDetails, err := ctrl.Client.ViewOrdersDetails(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if OrderDetails.UserID != user.UserId {
		log.Printf("order %s res %v", OrderDetails.OrderStatus, user.UserId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order is not of this user"})
		return
	}

	if strings.ToLower(order.OrderStatus) == "cancelled" {
		OrderDetails.OrderStatus = "cancelled"
		if err := ctrl.Client.ProcessOrder(*OrderDetails); err != nil {
			utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Post request failed", "Error", err.Error())
			return
		}
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Order Cancelled", "order", order)
		return
	}

	orderTransitions := payload.GetOrderTransitions()
	if OrderDetails.OrderStatus == "Delivered" {
		var driver model.User
		err := ctrl.Repo.GetUser("user_id", OrderDetails.DeliverDriverID, &driver)
		if err != nil {
			utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Delivery driver not found", "Error", err.Error())
			return
		}

		driver.RoleStatus = "available"
		ctrl.Repo.UpdateRoleStatus(&driver)
	}

	if newStatus, exists := orderTransitions[order.OrderStatus]; exists {
		OrderDetails.OrderStatus = newStatus
	}
	if err := ctrl.Client.ProcessOrder(*OrderDetails); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Patch request failed", "Error", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Order": order,
	})
}

func (ctrl *Controller) ProcessOrderDriver(c *gin.Context) {

	username, err := utils.VerificationUsername(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	activeRole, exists := c.Get("activeRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if activeRole != "Delivery driver" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient permission"})
		return
	}

	user := model.User{}
	err = ctrl.Repo.GetUser("username", username, &user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var order payload.ProcessOrder

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, "Error while binding order status")
		return
	}
	OrderDetails, err := ctrl.Client.ViewOrdersDetails(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if OrderDetails.DeliverDriverID != user.UserId {
		log.Printf("order %s res %v", OrderDetails.OrderStatus, user.UserId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order is not of this driver"})
		return
	}

	if strings.ToLower(order.OrderStatus) == "cancelled" {
		OrderDetails.OrderStatus = "cancelled"
		if err := ctrl.Client.ProcessOrder(*OrderDetails); err != nil {
			utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Post request failed", "Error", err.Error())
			return
		}
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Order Cancelled", "order", order)
		return
	}

	orderTransitions := payload.GetOrderTransitions()
	if OrderDetails.OrderStatus == "Delivered" {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Delivery driver can not complete the order", "", nil)
		return
	}

	if newStatus, exists := orderTransitions[order.OrderStatus]; exists {
		OrderDetails.OrderStatus = newStatus
	}
	if err := ctrl.Client.ProcessOrder(*OrderDetails); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Patch request failed", "Error", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Order": order,
	})
}
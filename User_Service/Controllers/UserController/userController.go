package UserControllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/E-Furqan/Food-Delivery-System/Client/AuthClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Repo        *database.Repository
	OrderClient *OrderClient.OrderClient
	AuthClient  *AuthClient.AuthClient
}

func NewController(repo *database.Repository, OrderClient *OrderClient.OrderClient, AuthClient *AuthClient.AuthClient) *Controller {
	return &Controller{
		Repo:        repo,
		OrderClient: OrderClient,
		AuthClient:  AuthClient,
	}
}

func (ctrl *Controller) Register(c *gin.Context) {

	var registrationData model.User

	if err := c.ShouldBindJSON(&registrationData); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", err.Error(), "", nil)
		return
	}

	if len(registrationData.Roles) > 0 && registrationData.ActiveRole == "" {
		var role model.Role
		if err := ctrl.Repo.GetRole(registrationData.Roles[0].RoleId, &role); err != nil {
			utils.GenerateResponse(http.StatusInternalServerError, c, "Error", "Role not found", "", nil)
			return
		}
		registrationData.ActiveRole = role.RoleType
		log.Print("active role set")
	}

	err := ctrl.Repo.CreateUser(&registrationData)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusCreated, registrationData)
}

func (ctrl *Controller) Login(c *gin.Context) {

	var input payload.Credentials
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", err.Error(), "", nil)
		return
	}

	var user model.User
	err := ctrl.Repo.GetUser("username", input.Username, &user)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "Invalid credentials", "", nil)
		return
	}

	if user.Password != input.Password {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "Invalid password", "", nil)
		return
	}
	var UserClaim payload.UserClaim
	UserClaim.Username = user.Username
	UserClaim.ActiveRole = user.ActiveRole
	UserClaim.ServiceType = "User"
	token, err := ctrl.AuthClient.GenerateToken(UserClaim)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", "Could not generate token", "", nil)
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
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "User not authenticated", "", nil)
		return
	}

	if activeRole != "Admin" {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", "You do not have the privileges to view users.", "", nil)
		return
	}

	var userData []model.User
	var OrderInfo payload.Order

	if err := c.ShouldBindJSON(&OrderInfo); err != nil {
		log.Print("binding error")
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", err.Error(), "", nil)
		return
	}

	userData, err := ctrl.Repo.PreloadInOrder(OrderInfo.ColumnName, OrderInfo.OrderType)

	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusOK, userData)
}

func (ctrl *Controller) UpdateUser(c *gin.Context) {

	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Error", err.Error(), "", nil)
		return
	}

	user := model.User{}
	err = ctrl.Repo.GetUser("user_id", UserId, &user)

	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Error", err.Error(), "", nil)
		return
	}

	var updateUserData model.User
	err = c.ShouldBindJSON(&updateUserData)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", err.Error(), "", nil)
		return
	}

	if err := ctrl.Repo.DeleteUserRoleInfo(user.UserId, "user_user_id"); err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", err.Error(), "", nil)
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
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", err.Error(), "", nil)
		return
	}

	for _, role := range updateUserData.Roles {
		if err := ctrl.Repo.AddUserRole(user.UserId, role.RoleId); err != nil {
			utils.GenerateResponse(http.StatusInternalServerError, c, "Error", err.Error(), "", nil)
			return
		}
	}
	c.JSON(http.StatusCreated, user)
}

func (ctrl *Controller) DeleteUser(c *gin.Context) {
	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Error", err.Error(), "", nil)
		return
	}

	user := model.User{}
	err = ctrl.Repo.GetUser("user_id", UserId, &user)

	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Error", err.Error(), "", nil)
		return
	}

	if err := ctrl.Repo.DeleteUserRoleInfo(user.UserId, "user_user_id"); err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", err.Error(), "", nil)
		return
	}

	if err := ctrl.Repo.DeleteUser(&user); err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (ctrl *Controller) Profile(c *gin.Context) {

	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Error", err.Error(), "", nil)
		return
	}

	var user model.User

	err = ctrl.Repo.GetUser("user_id", UserId, &user)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusFound, user)
}

func (ctrl *Controller) SearchForUser(c *gin.Context) {
	role, exists := c.Get("activeRole")
	if !exists {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "User not authenticated", "", nil)
		return
	}
	if role != "Admin" {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", "You do not have the privileges to Search for users.", "", nil)
		return
	}

	var input payload.UserSearch
	err := c.ShouldBindJSON(&input)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", err.Error(), "", nil)
		return
	}

	var user model.User
	err = ctrl.Repo.GetUser(input.ColumnName, input.SearchParameter, &user)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusFound, user)
}

func (ctrl *Controller) SwitchRole(c *gin.Context) {

	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Error", err.Error(), "", nil)
		return
	}

	var user model.User
	err = ctrl.Repo.GetUser("user_id", UserId, &user)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Error", err.Error(), "", nil)
		return
	}

	var RoleSwitch payload.RoleSwitch
	err = c.ShouldBindJSON(&RoleSwitch)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", err.Error(), "", nil)
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
		utils.GenerateResponse(http.StatusNotFound, c, "Error", "Role not found in user's roles", "", nil)
		return
	}

	user.ActiveRole = newRole.RoleType

	if err := ctrl.Repo.UpdateUserActiveRole(&user); err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", "Failed to update user active role", "", nil)
		return
	}

	var UserClaim payload.UserClaim
	UserClaim.Username = user.Username
	UserClaim.ActiveRole = user.ActiveRole
	UserClaim.ServiceType = "User"
	token, err := ctrl.AuthClient.GenerateToken(UserClaim)
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

func (ctrl *Controller) ViewUserOrders(c *gin.Context) {
	UserId, err := utils.VerifyUserId(c)

	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", err.Error(), "", nil)
		return
	}

	var User model.User
	err = ctrl.Repo.GetUser("user_id", UserId, &User)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Error", err.Error(), "", nil)
		return
	}
	var userId payload.ProcessOrder

	userId.UserID = User.UserId
	Orders, err := ctrl.OrderClient.ViewUserOrders(userId)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"User orders: ": Orders,
	})
}

func (ctrl *Controller) ProcessOrderUser(c *gin.Context) {
	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Error", err.Error(), "", nil)
		return
	}

	user := model.User{}

	err = ctrl.Repo.GetUser("user_id", UserId, &user)

	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Error", "User not found", "", nil)
		return
	}

	var order payload.ProcessOrder

	if err := c.ShouldBindJSON(&order); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", err.Error(), "", nil)
		return
	}
	OrderDetails, err := ctrl.OrderClient.ViewOrdersDetails(order)

	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", err.Error(), "", nil)
		return
	}
	if OrderDetails.UserID != user.UserId {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", "Order is not of this user", "", nil)
		return
	}

	if strings.ToLower(order.OrderStatus) == "cancelled" {
		OrderDetails.OrderStatus = "cancelled"
		if err := ctrl.OrderClient.ProcessOrder(*OrderDetails); err != nil {
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

	if newStatus, exists := orderTransitions[OrderDetails.OrderStatus]; exists {
		OrderDetails.OrderStatus = newStatus
	}
	if err := ctrl.OrderClient.ProcessOrder(*OrderDetails); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Patch request failed", "Error", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Order": OrderDetails,
	})
}

func (ctrl *Controller) ProcessOrderDriver(c *gin.Context) {

	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", err.Error(), "", nil)
		return
	}
	activeRole, exists := c.Get("activeRole")
	if !exists {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "User not authenticated", "", nil)
		return
	}

	if activeRole != "Delivery driver" {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", "insufficient permission", "", nil)
		return
	}

	user := model.User{}
	err = ctrl.Repo.GetUser("user_id", UserId, &user)

	if err != nil {

		utils.GenerateResponse(http.StatusNotFound, c, "Error", "User not found", "", nil)
		return
	}

	var order payload.ProcessOrder

	if err := c.ShouldBindJSON(&order); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", err.Error(), "", nil)
		return
	}

	OrderDetails, err := ctrl.OrderClient.ViewOrdersDetails(order)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", err.Error(), "", nil)
		return
	}
	if OrderDetails.DeliverDriverID == 0 {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", "Order have no delivery driver", "", nil)
		return
	}

	if OrderDetails.DeliverDriverID != user.UserId {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", "Order is not of this driver", "", nil)
		return
	}

	if strings.ToLower(order.OrderStatus) == "cancelled" {
		OrderDetails.OrderStatus = "cancelled"
		if err := ctrl.OrderClient.ProcessOrder(*OrderDetails); err != nil {
			utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Post request failed", "Error", err.Error())
			return
		}
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Order Cancelled", "order", order)
		return
	}

	orderTransitions := payload.GetOrderTransitions()
	if OrderDetails.OrderStatus == "Delivered" {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", "Delivery driver can not complete the order", "", nil)
		return
	}

	if newStatus, exists := orderTransitions[OrderDetails.OrderStatus]; exists {
		OrderDetails.OrderStatus = newStatus
	}
	if err := ctrl.OrderClient.ProcessOrder(*OrderDetails); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Patch request failed", "Error", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Order": OrderDetails,
	})
}

func (ctrl *Controller) ViewDriverOrders(c *gin.Context) {

	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "user not authenticated", "Error", err.Error())
		return
	}

	activeRole, exists := c.Get("activeRole")
	if !exists {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", "User role does not exist", "", nil)
		return
	}

	if activeRole != "Delivery driver" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "insufficient permission", "", nil)
		return
	}

	var User model.User
	err = ctrl.Repo.GetUser("user_id", UserId, &User)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "User does not exist", "", nil)
		return
	}
	var userId payload.ProcessOrder

	userId.UserID = User.UserId
	Orders, err := ctrl.OrderClient.ViewDriverOrders(userId)
	if err != nil {
		utils.GenerateResponse(http.StatusBadGateway, c, "Error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Driver orders: ": Orders,
	})
}

func (ctrl *Controller) ViewOrdersWithoutDriver(c *gin.Context) {
	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "user not authenticated", "Error", err.Error())
		return
	}

	activeRole, exists := c.Get("activeRole")
	if !exists {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", "User role does not exist", "", nil)
		return
	}

	if activeRole != "Delivery driver" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "insufficient permission", "", nil)
		return
	}

	var User model.User
	err = ctrl.Repo.GetUser("user_id", UserId, &User)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "User does not exist", "", nil)
		return
	}
	var userId payload.ProcessOrder
	Orders, err := ctrl.OrderClient.ViewOrdersWithoutRider(userId)
	if err != nil {
		utils.GenerateResponse(http.StatusBadGateway, c, "Error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Driver orders: ": Orders,
	})
}

func (ctrl *Controller) AssignDriver(c *gin.Context) {
	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "user not authenticated", "Error", err.Error())
		return
	}

	activeRole, exists := c.Get("activeRole")
	if !exists {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", "User role does not exist", "", nil)
		return
	}

	if activeRole != "Delivery driver" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "insufficient permission", "", nil)
		return
	}

	var driver model.User
	err = ctrl.Repo.GetUser("user_id", UserId, &driver)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "User does not exist", "", nil)
		return
	}
	var orderId payload.ProcessOrder

	if err := c.ShouldBindJSON(&orderId); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", err.Error(), "", nil)
		return
	}

	OrderDetails, err := ctrl.OrderClient.ViewOrdersDetails(orderId)
	if err != nil {
		utils.GenerateResponse(http.StatusBadGateway, c, "Error", err.Error(), "", nil)
		return
	}

	if OrderDetails.DeliverDriverID != 0 {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", "Order already have a delivery driver", "", nil)
		return
	}

	OrderDetails.DeliverDriverID = driver.UserId
	if err := ctrl.OrderClient.ProcessOrder(*OrderDetails); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Patch request failed", "Error", err.Error())
		return
	}

	driver.RoleStatus = "not available"
	ctrl.Repo.UpdateRoleStatus(&driver)

	c.JSON(http.StatusOK, gin.H{
		"Order": OrderDetails,
	})
}

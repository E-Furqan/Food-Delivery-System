package UserControllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/E-Furqan/Food-Delivery-System/Client/AuthClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
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
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}

	if len(registrationData.Roles) > 0 && registrationData.ActiveRole == "" {
		var role model.Role
		if err := ctrl.Repo.GetRole(registrationData.Roles[0].RoleId, &role); err != nil {
			utils.GenerateResponse(http.StatusInternalServerError, c, "error", "Role not found", "", nil)
			return
		}
		registrationData.ActiveRole = role.RoleType
		log.Print("active role set")
	}

	err := ctrl.Repo.CreateUser(&registrationData)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusCreated, registrationData)
}

func (ctrl *Controller) Login(c *gin.Context) {

	var input model.Credentials
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}

	var user model.User
	err := ctrl.Repo.GetUser("username", input.Username, &user)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "Invalid credentials", "", nil)
		return
	}
	if user.Password != input.Password {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "Invalid credentials", "", nil)
		return
	}

	UserClaim := utils.CreateUserClaim(user)

	token, err := ctrl.AuthClient.GenerateToken(UserClaim)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", "Could not generate token", "", nil)
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
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "User not authenticated", "", nil)
		return
	}

	if activeRole != "Admin" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "You do not have the privileges to view users.", "", nil)
		return
	}

	var userData []model.User
	var OrderInfo model.Order

	if err := c.ShouldBindJSON(&OrderInfo); err != nil {
		log.Print("binding error")
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}

	userData, err := ctrl.Repo.FetchUsersWithRoles(OrderInfo.ColumnName, OrderInfo.OrderType)

	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusOK, userData)
}

func (ctrl *Controller) UpdateUser(c *gin.Context) {

	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", err.Error(), "", nil)
		return
	}

	user := model.User{}
	err = ctrl.Repo.GetUser("user_id", UserId, &user)

	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "error", err.Error(), "", nil)
		return
	}

	var updateUserData model.User
	err = c.ShouldBindJSON(&updateUserData)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}

	if err := ctrl.Repo.DeleteUserRoleInfo(user.UserId, "user_user_id"); err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
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
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
		return
	}

	for _, role := range updateUserData.Roles {
		if err := ctrl.Repo.AddUserRole(user.UserId, role.RoleId); err != nil {
			utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
			return
		}
	}
	c.JSON(http.StatusCreated, user)
}

func (ctrl *Controller) DeleteUser(c *gin.Context) {
	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", err.Error(), "", nil)
		return
	}

	user := model.User{}
	err = ctrl.Repo.GetUser("user_id", UserId, &user)

	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "error", err.Error(), "", nil)
		return
	}

	if err := ctrl.Repo.DeleteUserRoleInfo(user.UserId, "user_user_id"); err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
		return
	}

	if err := ctrl.Repo.DeleteUser(&user); err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (ctrl *Controller) Profile(c *gin.Context) {

	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", err.Error(), "", nil)
		return
	}

	var user model.User

	err = ctrl.Repo.GetUser("user_id", UserId, &user)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusFound, user)
}

func (ctrl *Controller) SearchForUser(c *gin.Context) {
	role, exists := c.Get("activeRole")
	if !exists {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "User not authenticated", "", nil)
		return
	}
	if role != "Admin" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "You do not have the privileges to Search for users.", "", nil)
		return
	}

	var input model.UserSearch
	err := c.ShouldBindJSON(&input)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}

	var user model.User
	err = ctrl.Repo.GetUser(input.ColumnName, input.SearchParameter, &user)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusFound, user)
}

func (ctrl *Controller) ViewUserOrders(c *gin.Context) {
	UserId, err := utils.VerifyUserId(c)

	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", err.Error(), "", nil)
		return
	}

	var User model.User
	err = ctrl.Repo.GetUser("user_id", UserId, &User)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "error", err.Error(), "", nil)
		return
	}
	var userId model.ProcessOrder

	userId.UserID = User.UserId
	Orders, err := ctrl.OrderClient.ViewUserOrders(userId)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"User orders: ": Orders,
	})
}

func (ctrl *Controller) ProcessOrderUser(c *gin.Context) {
	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", err.Error(), "", nil)
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "authorization token not provided", "error", nil)
		return
	}
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "invalid authorization header format", "error", nil)
		return
	}
	token := tokenParts[1]

	user := model.User{}

	err = ctrl.Repo.GetUser("user_id", UserId, &user)

	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "error", "User not found", "", nil)
		return
	}

	var order model.ProcessOrder

	if err := c.ShouldBindJSON(&order); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}
	OrderDetails, err := ctrl.OrderClient.ViewOrdersDetails(order, token)

	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
		return
	}
	if OrderDetails.UserID != user.UserId {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", "Order is not of this user", "", nil)
		return
	}

	if strings.ToLower(order.OrderStatus) == "cancelled" {
		OrderDetails.OrderStatus = "cancelled"
		if err := ctrl.OrderClient.ProcessOrder(*OrderDetails); err != nil {
			utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Post request failed", "error", err.Error())
			return
		}
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Order Cancelled", "order", order)
		return
	}

	orderTransitions := model.GetOrderTransitions()
	if OrderDetails.OrderStatus == "Delivered" {
		var driver model.User
		err := ctrl.Repo.GetUser("user_id", OrderDetails.DeliverDriverID, &driver)
		if err != nil {
			utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Delivery driver not found", "error", err.Error())
			return
		}

		driver.RoleStatus = "available"
		ctrl.Repo.UpdateRoleStatus(&driver)
	}

	if newStatus, exists := orderTransitions[OrderDetails.OrderStatus]; exists {
		OrderDetails.OrderStatus = newStatus
	}
	if err := ctrl.OrderClient.ProcessOrder(*OrderDetails); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Patch request failed", "error", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Order": OrderDetails,
	})
}

func (ctrl *Controller) ProcessOrderDriver(c *gin.Context) {

	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", err.Error(), "", nil)
		return
	}
	activeRole, exists := c.Get("activeRole")
	if !exists {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "User not authenticated", "", nil)
		return
	}

	if activeRole != "Delivery driver" {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", "insufficient permission", "", nil)
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "authorization token not provided", "error", nil)
		return
	}
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "invalid authorization header format", "error", nil)
		return
	}
	token := tokenParts[1]

	user := model.User{}
	err = ctrl.Repo.GetUser("user_id", UserId, &user)

	if err != nil {

		utils.GenerateResponse(http.StatusNotFound, c, "error", "User not found", "", nil)
		return
	}

	var order model.ProcessOrder

	if err := c.ShouldBindJSON(&order); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}

	OrderDetails, err := ctrl.OrderClient.ViewOrdersDetails(order, token)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
		return
	}
	if OrderDetails.DeliverDriverID == 0 {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", "Order have no delivery driver", "", nil)
		return
	}

	if OrderDetails.DeliverDriverID != user.UserId {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", "Order is not of this driver", "", nil)
		return
	}

	if strings.ToLower(order.OrderStatus) == "cancelled" {
		OrderDetails.OrderStatus = "cancelled"
		if err := ctrl.OrderClient.ProcessOrder(*OrderDetails); err != nil {
			utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Post request failed", "error", err.Error())
			return
		}
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Order Cancelled", "order", order)
		return
	}

	orderTransitions := model.GetOrderTransitions()
	if OrderDetails.OrderStatus == "Delivered" {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", "Delivery driver can not complete the order", "", nil)
		return
	}

	if newStatus, exists := orderTransitions[OrderDetails.OrderStatus]; exists {
		OrderDetails.OrderStatus = newStatus
	}
	if err := ctrl.OrderClient.ProcessOrder(*OrderDetails); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Patch request failed", "error", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Order": OrderDetails,
	})
}

func (ctrl *Controller) ViewDriverOrders(c *gin.Context) {

	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "user not authenticated", "error", err.Error())
		return
	}

	activeRole, exists := c.Get("activeRole")
	if !exists {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "User role does not exist", "", nil)
		return
	}

	if activeRole != "Delivery driver" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "insufficient permission", "", nil)
		return
	}

	var User model.User
	err = ctrl.Repo.GetUser("user_id", UserId, &User)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "User does not exist", "", nil)
		return
	}
	var userId model.ProcessOrder

	userId.DeliverDriverID = User.UserId
	Orders, err := ctrl.OrderClient.ViewDriverOrders(userId)
	if err != nil {
		utils.GenerateResponse(http.StatusBadGateway, c, "error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Driver orders: ": Orders,
	})
}

func (ctrl *Controller) ViewOrdersWithoutDriver(c *gin.Context) {
	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "user not authenticated", "error", err.Error())
		return
	}

	activeRole, exists := c.Get("activeRole")
	if !exists {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", "User role does not exist", "", nil)
		return
	}

	if activeRole != "Delivery driver" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "insufficient permission", "", nil)
		return
	}

	var User model.User
	err = ctrl.Repo.GetUser("user_id", UserId, &User)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "User does not exist", "", nil)
		return
	}
	var userId model.ProcessOrder
	Orders, err := ctrl.OrderClient.ViewOrdersWithoutRider(userId)
	if err != nil {
		utils.GenerateResponse(http.StatusBadGateway, c, "error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Driver orders: ": Orders,
	})
}

func (ctrl *Controller) AssignDriver(c *gin.Context) {
	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "user not authenticated", "error", err.Error())
		return
	}

	activeRole, exists := c.Get("activeRole")
	if !exists {
		utils.GenerateResponse(http.StatusNotFound, c, "error", "User role does not exist", "", nil)
		return
	}

	if activeRole != "Delivery driver" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "insufficient permission", "", nil)
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "authorization token not provided", "error", nil)
		return
	}
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "invalid authorization header format", "error", nil)
		return
	}
	token := tokenParts[1]

	var driver model.User
	err = ctrl.Repo.GetUser("user_id", UserId, &driver)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "User does not exist", "", nil)
		return
	}
	var orderId model.ProcessOrder

	if err := c.ShouldBindJSON(&orderId); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}

	OrderDetails, err := ctrl.OrderClient.ViewOrdersDetails(orderId, token)
	if err != nil {
		utils.GenerateResponse(http.StatusBadGateway, c, "error", err.Error(), "", nil)
		return
	}

	if OrderDetails.DeliverDriverID != 0 {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", "Order already have a delivery driver", "", nil)
		return
	}

	OrderDetails.DeliverDriverID = driver.UserId
	if err := ctrl.OrderClient.ProcessOrder(*OrderDetails); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Patch request failed", "error", err.Error())
		return
	}

	driver.RoleStatus = "not available"
	ctrl.Repo.UpdateRoleStatus(&driver)

	c.JSON(http.StatusOK, gin.H{
		"Order": OrderDetails,
	})
}

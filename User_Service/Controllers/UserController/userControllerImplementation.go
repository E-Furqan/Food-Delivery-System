package UserControllers

import (
	"context"
	"log"
	"net/http"
	"strings"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

func (ctrl *Controller) RegisterWorkflow(c *gin.Context) {
	var registrationData model.User

	if err := c.ShouldBindJSON(&registrationData); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}
	log.Print("register data:", registrationData)

	client_var, err := client.Dial(client.Options{})
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "message", "unable to create Temporal client", "error", err)
		return
	}
	if client_var == nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "message", "Temporal client is nil", "error", err)
		return
	}

	options := client.StartWorkflowOptions{
		ID:                    "registration-workflow " + registrationData.Username,
		TaskQueue:             model.RegisterTaskQueue,
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
	}

	_, err = client_var.ExecuteWorkflow(context.Background(), options, ctrl.WorkFlows.RegisterWorkflow, registrationData)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "message", "error in workflow", "error", err)
		return
	}
	log.Print("user created")

	c.JSON(http.StatusOK, "user registered")
}

func (ctrl *Controller) Register(c *gin.Context) {

	var registrationData model.User

	if err := c.ShouldBindJSON(&registrationData); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}
	log.Print(registrationData)
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

	_, err := utils.VerifyActiveAdminRole(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "user not authenticated", "error", err.Error())
		return
	}

	var userData []model.User
	var OrderInfo model.Order

	if err := c.ShouldBindJSON(&OrderInfo); err != nil {
		log.Print("binding error")
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}

	userData, err = ctrl.Repo.FetchUsersWithRoles(OrderInfo.ColumnName, OrderInfo.OrderType)

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
	_, err := utils.VerifyActiveAdminRole(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "user not authenticated", "error", err.Error())
		return
	}

	var input model.UserSearch
	err = c.ShouldBindJSON(&input)
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
	var userId model.UpdateOrder

	userId.UserID = User.UserId
	Orders, err := ctrl.OrderClient.ViewOrders(userId, c)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"User orders: ": Orders,
	})
}

func (ctrl *Controller) UpdateOrderStatus(c *gin.Context) {
	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", err.Error(), "", nil)
		return
	}

	user := model.User{}
	err = ctrl.Repo.GetUser("user_id", UserId, &user)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "error", "User not found", "", nil)
		return
	}

	var order model.UpdateOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}
	updatedOrder, err := ctrl.OrderClient.UpdateOrderStatus(order, c)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Patch request failed", "error", err.Error())
		return
	}

	if strings.ToLower(updatedOrder.OrderStatus) == "completed" {
		if updatedOrder.DeliverDriverID == user.UserId {
			user.RoleStatus = "Available"
			ctrl.Repo.UpdateRoleStatus(&user)
		} else {
			user := model.User{}
			err = ctrl.Repo.GetUser("user_id", updatedOrder.DeliverDriverID, &user)
			if err != nil {
				utils.GenerateResponse(http.StatusNotFound, c, "error", "User not found", "", nil)
				return
			}
			user.RoleStatus = "Available"
			ctrl.Repo.UpdateRoleStatus(&user)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "order status updated",
	})
}

func (ctrl *Controller) ViewDriverOrders(c *gin.Context) {

	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "user not authenticated", "error", err.Error())
		return
	}

	activeRole, err := utils.FetchActiveRole(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", err.Error(), "", nil)
		return
	}

	err = utils.VerifyIfDriver(activeRole)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "insufficient permission", "", nil)
		return
	}

	var User model.User
	err = ctrl.Repo.GetUser("user_id", UserId, &User)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "User does not exist", "", nil)
		return
	}

	var userId model.UpdateOrder
	userId.DeliverDriverID = User.UserId
	Orders, err := ctrl.OrderClient.ViewOrders(userId, c)
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

	activeRole, err := utils.FetchActiveRole(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", err.Error(), "", nil)
		return
	}

	err = utils.VerifyIfDriver(activeRole)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "insufficient permission", "", nil)
		return
	}

	var User model.User
	err = ctrl.Repo.GetUser("user_id", UserId, &User)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "User does not exist", "", nil)
		return
	}

	var userId model.UpdateOrder
	Orders, err := ctrl.OrderClient.ViewOrdersWithoutDriver(userId, c)
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

	activeRole, err := utils.FetchActiveRole(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", err.Error(), "", nil)
		return
	}

	err = utils.VerifyIfDriver(activeRole)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "insufficient permission", "", nil)
		return
	}

	var driver model.User
	err = ctrl.Repo.GetUser("user_id", UserId, &driver)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "User does not exist", "", nil)
		return
	}

	var orderId model.UpdateOrder
	if err := c.ShouldBindJSON(&orderId); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}

	if err := ctrl.OrderClient.AssignDriver(orderId, c); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Patch request failed", "error", err.Error())
		return
	}

	driver.RoleStatus = "not available"
	ctrl.Repo.UpdateRoleStatus(&driver)

	c.JSON(http.StatusOK, gin.H{
		"Message": "Delivery driver assigned to the order",
	})
}

func (ctrl *Controller) FetchActiveUser(c *gin.Context) {

	var RoleFilter model.UserRoleFilter
	if err := c.ShouldBindJSON(&RoleFilter); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}
	result, err := utils.FetchActiveUserCountHelper(RoleFilter, ctrl.Repo)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusOK, result)
}

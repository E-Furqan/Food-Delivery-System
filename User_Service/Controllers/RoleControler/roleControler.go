package roleController

import (
	"log"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	Repo *database.Repository
}

func NewController(repo *database.Repository) *RoleController {
	return &RoleController{Repo: repo}
}

func (rCtrl *RoleController) AddRolesByAdmin(c *gin.Context) {

	activeRole, exists := c.Get("activeRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if activeRole != "Admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You do not have the privileges to add new roles."})
		return
	}

	var input payload.Role
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var role model.Role

	role.RoleId = input.RoleId
	role.RoleType = input.RoleType

	err := rCtrl.Repo.CreateRole(&role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role added successfully"})

}

func (rCtrl *RoleController) GetRole(c *gin.Context) {

	activeRole, exists := c.Get("activeRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if activeRole != "Admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You do not have the privileges to view roles."})
		return
	}

	var OrderInfo payload.Order
	if err := c.ShouldBindJSON(&OrderInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	RoleData, err := rCtrl.Repo.RoleInOrder(OrderInfo.ColumnName, OrderInfo.OrderType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, RoleData)
}

func (rCtrl *RoleController) DeleteRole(c *gin.Context) {

	activeRole, exists := c.Get("activeRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if activeRole != "Admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You do not have the privileges to Delete roles."})
		return
	}

	var input payload.Role
	var Role model.Role
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = rCtrl.Repo.GetRole(input.RoleId, &Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role does not exist"})
		return
	}

	if err := rCtrl.Repo.DeleteUserRoleInfo(input.RoleId, "role_role_id"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := rCtrl.Repo.DeleteRole(&Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})

}

func (rCtrl *RoleController) AddDefaultRoles(c *gin.Context) {
	var roles []model.Role

	for _, RolesFromPayLoad := range payload.RolesList {
		var existingRole model.Role
		err := rCtrl.Repo.GetRole(RolesFromPayLoad.RoleId, &existingRole)
		if err == nil {

			log.Printf("Role %v already exists, skipping.", RolesFromPayLoad.RoleId)
			continue
		}

		roles = append(roles, model.Role{
			RoleId:   RolesFromPayLoad.RoleId,
			RoleType: RolesFromPayLoad.RoleType,
		})
	}

	if len(roles) == 0 {
		log.Println("No new roles to add, exiting function.")
		return
	}

	err := rCtrl.Repo.BulkCreateRoles(roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add default roles"})
		return
	}

	log.Printf("Message: Default roles added successfully ")
}

// func (rCtrl *RoleController) AddRoleToUser(c *gin.Context) {
// 	username, exists := c.Get("username")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
// 		return
// 	}

// 	usernameStr, ok := username.(string)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username type"})
// 		return
// 	}

// 	var user model.User
// 	err := rCtrl.Repo.FindUser("username", usernameStr, &user)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User not found %v %s", err, usernameStr)})
// 		retu/ func (rCtrl *RoleController) AddRoleToUser(c *gin.Context) {
// 	username, exists := c.Get("username")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
// 		return
// 	}

// 	usernameStr, ok := username.(string)
// 	if !ok {
// 		c.JS
// 	for _, roleid := range updateData.Roles {
// 		var role model.Role
// 		if err := rCtrl.Repo.FindRole(roleid.RoleId, &role); err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
// 			return
// 		}

// 		existingRole := model.UserRole{}
// 		if err := rCtrl.Repo.DB.Where("user_id = ? AND role_id = ?", user.UserId, roleid.RoleId).First(&existingRole).Error; err == nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "User already has this role"})
// 			return
// 		}

// 		userRole := model.UserRole{
// 			UserId: user.UserId,
// 			RoleId: roleid.RoleId,
// 		}

// 		if err := rCtrl.Repo.DB.Create(&userRole).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add role to user"})
// 			return
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Role added to user successfully"})
// }

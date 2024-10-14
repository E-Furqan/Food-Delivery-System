package roleController

import (
	"net/http"

	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/enviorment_variable"
	model "github.com/E-Furqan/Food-Delivery-System/models"
	"github.com/E-Furqan/Food-Delivery-System/payload"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	Repo *database.Repository
}

func NewController(repo *database.Repository) *RoleController {
	return &RoleController{Repo: repo}
}

func (rCtrl *RoleController) AddRoles(RoleId string, c *gin.Context, role *model.Role) bool {

	err := rCtrl.Repo.FindRole(RoleId, role)
	if err != nil {

		if RoleId == "1" {

			role.RoleId = RoleId
			role.RoleType = "Customer"

			rCtrl.Repo.CreateRole(role)

		} else if RoleId == "2" {

			role.RoleId = RoleId
			role.RoleType = "Delivery driver"

			rCtrl.Repo.CreateRole(role)

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
			return false
		}

	}
	return true
}

func (rCtrl *RoleController) GetRole(c *gin.Context) {

	var user_data []model.Role
	var OrderInfo payload.Order

	if err := c.ShouldBindJSON(&OrderInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_data, err := rCtrl.Repo.RoleInOrder(OrderInfo.ColumnName, OrderInfo.OrderType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user_data)
}

// role con
func (rCtrl *RoleController) DeleteRole(c *gin.Context) {

	Admin := environmentVariable.Get_env("ADMIN")
	Admin_password := environmentVariable.Get_env("ADMIN_PASS")

	var input payload.Input

	// Bind incoming JSON to input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var role model.Role
	if input.Username == Admin && input.Password == Admin_password {
		// Fetch the role by role id
		err := rCtrl.Repo.FindRole(input.RoleId, &role)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		// Delete the role
		if err := rCtrl.Repo.DeleteRole(&role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the role"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
	}

}

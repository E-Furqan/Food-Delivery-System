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

func (rCtrl *RoleController) CheckIfRoleExist(Role_id string, c *gin.Context, role *model.Role) bool {

	err := rCtrl.Repo.FindRole(Role_id, role)
	if err != nil {

		if Role_id == "1" {

			role.Role_id = Role_id
			role.Role_type = "Customer"

			rCtrl.Repo.CreateRole(role)

		} else if Role_id == "2" {

			role.Role_id = Role_id
			role.Role_type = "Delivery driver"

			rCtrl.Repo.CreateRole(role)

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
			return false
		}

	}
	return true
}

// siftcres acs
func (rCtrl *RoleController) GetRole(c *gin.Context) {
	var user_data []model.Role
	user_data, err := rCtrl.Repo.RoleInAscOrder()
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
		err := rCtrl.Repo.FindRole(input.Role_id, &role)
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

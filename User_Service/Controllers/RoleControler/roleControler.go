package roleController

import (
	"log"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
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
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "User not authenticated", "", nil)
		return
	}

	if activeRole != "Admin" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "You do not have the privileges to add new roles.", "", nil)
		return
	}

	var input payload.Role
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", err.Error(), "", nil)
		return
	}
	var role model.Role

	role.RoleId = input.RoleId
	role.RoleType = input.RoleType

	err := rCtrl.Repo.CreateRole(&role)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", err.Error(), "", nil)
		return
	}

	utils.GenerateResponse(http.StatusOK, c, "Message", "Role added successfully", "", nil)

}

func (rCtrl *RoleController) GetRole(c *gin.Context) {

	activeRole, exists := c.Get("activeRole")
	if !exists {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "User not authenticated", "", nil)
		return
	}

	if activeRole != "Admin" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "You do not have the privileges to view roles.", "", nil)
		return
	}

	var OrderInfo payload.Order
	if err := c.ShouldBindJSON(&OrderInfo); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", err.Error(), "", nil)
		return
	}

	RoleData, err := rCtrl.Repo.RoleInOrder(OrderInfo.ColumnName, OrderInfo.OrderType)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusOK, RoleData)
}

func (rCtrl *RoleController) DeleteRole(c *gin.Context) {

	activeRole, exists := c.Get("activeRole")
	if !exists {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "User not authenticated", "", nil)
		return
	}

	if activeRole != "Admin" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Error", "You do not have the privileges to delete roles.", "", nil)
		return
	}

	var input payload.Role
	var Role model.Role
	err := c.ShouldBindJSON(&input)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Error", err.Error(), "", nil)
		return
	}

	err = rCtrl.Repo.GetRole(input.RoleId, &Role)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "Message", "Role does not exist", "Error", err.Error())
		return
	}

	if err := rCtrl.Repo.DeleteUserRoleInfo(input.RoleId, "role_role_id"); err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Error", err.Error(), "", nil)
		return
	}

	if err := rCtrl.Repo.DeleteRole(&Role); err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Message", "Failed to delete the role", "Error", err.Error())
		return
	}
	utils.GenerateResponse(http.StatusOK, c, "Message", "Role deleted successfully", "", nil)

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
		utils.GenerateResponse(http.StatusInternalServerError, c, "Message", "Failed to add default roles", "Error", err.Error())
		return
	}

	log.Printf("Message: Default roles added successfully ")
}

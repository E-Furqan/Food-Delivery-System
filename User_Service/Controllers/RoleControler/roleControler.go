package roleController

import (
	"log"
	"net/http"

	"github.com/E-Furqan/Food-Delivery-System/Client/AuthClient"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	Repo       *database.Repository
	AuthClient *AuthClient.AuthClient
}

func NewController(repo *database.Repository, AuthClient *AuthClient.AuthClient) *RoleController {
	return &RoleController{
		Repo:       repo,
		AuthClient: AuthClient}
}

func (rCtrl *RoleController) AddRolesByAdmin(c *gin.Context) {

	_, err := utils.VerifyActiveRole(c)
	if err != nil {
		utils.GenerateResponse(http.StatusUnauthorized, c, "Message", "user not authenticated", "error", err.Error())
		return
	}

	var input model.Role
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}
	var role model.Role

	role.RoleId = input.RoleId
	role.RoleType = input.RoleType

	err = rCtrl.Repo.CreateRole(&role)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
		return
	}

	utils.GenerateResponse(http.StatusOK, c, "Message", "Role added successfully", "", nil)
}

func (rCtrl *RoleController) GetRoles(c *gin.Context) {

	activeRole, exists := c.Get("activeRole")
	if !exists {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "User not authenticated", "", nil)
		return
	}
	log.Print(activeRole)
	if activeRole != "Admin" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "You do not have the privileges to view roles.", "", nil)
		return
	}

	RoleData, err := rCtrl.Repo.GetAllRoles()
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
		return
	}

	c.JSON(http.StatusOK, RoleData)
}

func (rCtrl *RoleController) DeleteRole(c *gin.Context) {

	activeRole, exists := c.Get("activeRole")
	if !exists {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "User not authenticated", "", nil)
		return
	}

	if activeRole != "Admin" {
		utils.GenerateResponse(http.StatusUnauthorized, c, "error", "You do not have the privileges to delete roles.", "", nil)
		return
	}

	var input model.Role
	var Role model.Role
	err := c.ShouldBindJSON(&input)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
		return
	}

	err = rCtrl.Repo.GetRole(input.RoleId, &Role)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "Message", "Role does not exist", "error", err.Error())
		return
	}

	if err := rCtrl.Repo.DeleteUserRoleInfo(input.RoleId, "role_role_id"); err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
		return
	}

	if err := rCtrl.Repo.DeleteRole(&Role); err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "Message", "Failed to delete the role", "error", err.Error())
		return
	}
	utils.GenerateResponse(http.StatusOK, c, "Message", "Role deleted successfully", "", nil)

}

func (rCtrl *RoleController) AddDefaultRoles(c *gin.Context) {
	var roles []model.Role

	for _, RolesFromPayLoad := range model.RolesList {
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
		utils.GenerateResponse(http.StatusInternalServerError, c, "Message", "Failed to add default roles", "error", err.Error())
		return
	}

	log.Printf("Message: Default roles added successfully ")
}

func (RoleController *RoleController) SwitchRole(c *gin.Context) {

	UserId, err := utils.VerifyUserId(c)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "error", err.Error(), "", nil)
		return
	}
	log.Print(UserId)
	var user model.User
	err = RoleController.Repo.GetUser("user_id", UserId, &user)
	if err != nil {
		utils.GenerateResponse(http.StatusNotFound, c, "error", err.Error(), "", nil)
		return
	}

	var RoleSwitch model.RoleSwitch
	err = c.ShouldBindJSON(&RoleSwitch)
	if err != nil {
		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
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
		utils.GenerateResponse(http.StatusNotFound, c, "error", "Role not found in user's roles", "", nil)
		return
	}

	user.ActiveRole = newRole.RoleType

	if err := RoleController.Repo.UpdateUserActiveRole(&user); err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, c, "error", "Failed to update user active role", "", nil)
		return
	}

	UserClaim := utils.CreateUserClaim(user)

	token, err := RoleController.AuthClient.GenerateToken(UserClaim)
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

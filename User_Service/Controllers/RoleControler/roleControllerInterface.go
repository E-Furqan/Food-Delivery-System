package RoleController

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/AuthClient"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	Repo       database.RepositoryInterface
	AuthClient AuthClient.AuthClientInterface
}

func NewController(repo database.RepositoryInterface, AuthClient AuthClient.AuthClientInterface) *RoleController {
	return &RoleController{
		Repo:       repo,
		AuthClient: AuthClient}
}

type RoleControllerInterface interface {
	AddRolesByAdmin(c *gin.Context)
	GetRoles(c *gin.Context)
	DeleteRole(c *gin.Context)
	AddDefaultRoles(c *gin.Context)
	SwitchRole(c *gin.Context)
}

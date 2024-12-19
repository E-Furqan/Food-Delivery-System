package dataController

import (
	driveClient "github.com/E-Furqan/Food-Delivery-System/Client/DriveClient"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Repo        database.RepositoryInterface
	DriveClient driveClient.DriveControllerInterface
}

func NewController(repo database.RepositoryInterface, driveClient driveClient.DriveControllerInterface) *Controller {
	return &Controller{
		Repo:        repo,
		DriveClient: driveClient,
	}
}

type DataControllerInterface interface {
	SaveConfiguration(ctx *gin.Context)
}

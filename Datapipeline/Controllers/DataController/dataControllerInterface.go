package dataController

import (
	driveClient "github.com/E-Furqan/Food-Delivery-System/Client/DriveClient"
	workflowClient "github.com/E-Furqan/Food-Delivery-System/Client/WorkFlowClient"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Repo        database.RepositoryInterface
	DriveClient driveClient.DriveClientInterface
	WorkFlow    workflowClient.RestaurantClientInterface
}

func NewController(repo database.RepositoryInterface, driveClient driveClient.DriveClientInterface,
	workFlow workflowClient.RestaurantClientInterface) *Controller {
	return &Controller{
		Repo:        repo,
		DriveClient: driveClient,
		WorkFlow:    workFlow,
	}
}

type DataControllerInterface interface {
	CreateSourceConfiguration(ctx *gin.Context)
	CreateDestinationConfiguration(ctx *gin.Context)
	CreatePipeline(ctx *gin.Context)
	StartDatapipelineSync(ctx *gin.Context)
	FetchDestinationConfiguration(ctx *gin.Context)
	FetchSourceConfiguration(ctx *gin.Context)
}

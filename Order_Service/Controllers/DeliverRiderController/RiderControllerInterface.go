package RiderController

import (
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
)

type RiderController struct {
	Repo database.RepositoryInterface
}

func NewController(repo database.RepositoryInterface) *RiderController {
	return &RiderController{
		Repo: repo,
	}
}

type RiderControllerInterface interface {
	AssignDeliveryDriver(c *gin.Context)
	FetchCompletedDeliversRider(c *gin.Context)
}

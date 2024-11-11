package ItemController

import (
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
)

type ItemController struct {
	Repo *database.Repository
}

func NewController(repo *database.Repository) *ItemController {
	return &ItemController{Repo: repo}
}

type ItemControllerInterface interface {
	AddItemsInMenu(c *gin.Context)
	DeleteItemsFromMenu(c *gin.Context)
}

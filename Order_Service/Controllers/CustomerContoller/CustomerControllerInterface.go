package CustomerController

import (
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	Repo database.RepositoryInterface
}

func NewController(repo database.RepositoryInterface) *CustomerController {
	return &CustomerController{
		Repo: repo,
	}
}

type CustomerControllerInterface interface {
	FetchTopFiveCustomers(c *gin.Context)
	FetchCustomerOrdersDetails(c *gin.Context)
}

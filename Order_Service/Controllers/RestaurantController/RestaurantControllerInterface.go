package RestaurantController

import (
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
)

type RestaurantController struct {
	Repo database.RepositoryInterface
}

func NewController(repo database.RepositoryInterface) *RestaurantController {
	return &RestaurantController{
		Repo: repo,
	}
}

type RestaurantControllerInterface interface {
	FetchTopPurchasedItems(c *gin.Context)
	FetchCompletedOrdersCountByRestaurant(c *gin.Context)
	FetchRevenueOfRestaurants(c *gin.Context)
}

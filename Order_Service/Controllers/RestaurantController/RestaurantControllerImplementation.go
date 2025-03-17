package RestaurantController

import (
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

// FetchTopPurchasedItems godoc
// @Summary Fetch Top Purchased Items
// @Description Retrieves the top purchased items across all orders with their purchase counts.
// @Tags Order Service
// @Accept json
// @Produce json
// @Success 200 {array} model.TopPurchasedItem "List of top purchased items with their counts"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /order/top-purchased-items [get]
func (riderCtrl *RestaurantController) FetchTopPurchasedItems(c *gin.Context) {

	result, err := riderCtrl.Repo.FetchTopPurchasedItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// FetchCompletedOrdersCountByRestaurant godoc
// @Summary Fetch Completed Orders Count by Restaurant
// @Description Retrieves the number of completed orders for each restaurant within a given time range.
// @Tags Order Service
// @Accept json
// @Produce json
// @Param TimeRange body model.TimeRange true "Time Range for filtering completed orders"
// @Success 200 {array} model.CompletedOrdersCount "List of restaurants with completed orders count"
// @Failure 400 {object} map[string]interface{} "Invalid Time Range Input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /order/completed-orders-by-restaurant [get]
func (riderCtrl *RestaurantController) FetchCompletedOrdersCountByRestaurant(c *gin.Context) {

	activeRoleStr, err := utils.FetchRoleFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	isAdminOrRestaurant := utils.IsRestaurantOrAdminRole(activeRoleStr)
	if !isAdminOrRestaurant {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can get completed delivers by riders analysis"})
		return
	}

	var TimeRange model.TimeRange
	if err := c.ShouldBindJSON(&TimeRange); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	result, err := riderCtrl.Repo.FetchCompletedOrdersCountByRestaurant(TimeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// FetchRevenueOfRestaurants godoc
// @Summary Fetch revenue data of restaurants
// @Description This API endpoint returns the revenue data of restaurants, available only to users with admin or restaurant roles
// @Tags restaurants
// @Accept json
// @Produce json
// @Success 200 {array} model.RestaurantRevenue "Revenue data of restaurants"
// @Failure 401 {object} gin.H "Unauthorized - only admins or restaurant users can access"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /restaurants/revenue [get]
func (riderCtrl *RestaurantController) FetchRevenueOfRestaurants(c *gin.Context) {
	activeRoleStr, err := utils.FetchRoleFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	isAdminOrRestaurant := utils.IsRestaurantOrAdminRole(activeRoleStr)
	if !isAdminOrRestaurant {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can get completed delivers by riders analysis"})
		return
	}

	result, err := riderCtrl.Repo.FetchRestaurantsRevenue()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

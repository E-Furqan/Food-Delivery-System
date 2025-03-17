package CustomerController

import (
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

// FetchCustomerOrdersDetails godoc
// @Summary Fetch orders details for a specific customer
// @Description Retrieves the details of orders placed by a specific customer, including item information
// @Tags Order Service
// @Accept json
// @Produce json
// @Param pageNumber body model.PageNumber true "Pagination information"
// @Success 200 {array} model.OrderDetails "List of orders placed by the customer with item details"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid page number or limit"
// @Failure 401 {object} map[string]interface{} "Unauthorized: User not authorized to access the resource"
// @Failure 404 {object} map[string]interface{} "Not Found: Customer not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /order/customer/orders [get]
// @Security ApiKeyAuth
func (cusCtrl *CustomerController) FetchCustomerOrdersDetails(c *gin.Context) {

	activeRoleStr, err := utils.FetchRoleFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	isAdmin := utils.IsCustomerOrAdminRole(activeRoleStr)
	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can get completed delivers by riders analysis"})
		return
	}

	ID, err := utils.FetchIDFromClaim(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, "ClaimId is not a valid uint")
		return
	}

	var PageNumber model.PageNumber
	if err := c.ShouldBindJSON(&PageNumber); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	offset := (PageNumber.PageNumber - 1) * PageNumber.Limit
	result, err := cusCtrl.Repo.FetchUserOrdersWithItemDetails(int(ID), PageNumber.Limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// FetchTopFiveCustomers godoc
// @Summary Fetch the top five customers
// @Description This API endpoint returns the top five customers based on the completion of their orders
// @Tags customers
// @Accept json
// @Produce json
// @Success 200 {array} model.Customer "Top five customers"
// @Failure 401 {object} gin.H "Unauthorized - only admins can access this data"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /customers/top-five [get]
func (cusCtrl *CustomerController) FetchTopFiveCustomers(c *gin.Context) {
	activeRoleStr, err := utils.FetchRoleFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	isAdmin := utils.IsAdminRole(activeRoleStr)
	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can get completed delivers by riders analysis"})
		return
	}

	result, err := cusCtrl.Repo.FetchTopFiveCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

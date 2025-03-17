package RiderController

import (
	"net/http"
	"strings"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

// AssignDeliveryDriver godoc
// @Summary Assign a delivery driver to an order
// @Description Assigns a delivery driver to an order if the order doesn't already have a driver, and the role of the user is "delivery driver"
// @Tags Order Service
// @Accept  json
// @Produce  json
// @Param assignDeliveryDriverRequest body model.AssignDeliveryDriver true "Assign delivery driver request"
// @Success 200 {object} model.Order "Assigned delivery driver to the order"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Router /order/assign/diver [patch]
func (restCtrl *RiderController) AssignDeliveryDriver(c *gin.Context) {
	Id, exists := c.Get("ClaimId")
	if !exists {
		c.JSON(http.StatusBadRequest, "userId id does not exist")
		return
	}
	IDint := Id.(uint)

	activeRoleStr, err := utils.FetchRoleFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var request model.AssignDeliveryDriver
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.ToLower(activeRoleStr) != "delivery driver" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	var order model.Order

	err = restCtrl.Repo.GetOrder(&order, request.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.DeliveryDriver != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order already have a driver"})
		return
	}
	order.DeliveryDriver = IDint
	if err := restCtrl.Repo.Update(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// FetchCompletedDeliversRider godoc
// @Summary Fetch completed deliveries by riders
// @Description Retrieves the number of completed deliveries for each rider. This endpoint is restricted to admins only.
// @Tags Order Service
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.CompletedDelivers "List of completed deliveries by riders"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /order/completed/delivers/rider [get]
// @Summary Fetch completed deliveries by riders
// @Description Retrieves the number of completed deliveries for each rider
// @Tags Order Service
// @Accept  json
// @Produce  json
// @Success 200 {array} model.CompletedDelivers "List of completed deliveries by riders"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /order/completed/delivers/rider [get]
// @Security ApiKeyAuth
func (restCtrl *RiderController) FetchCompletedDeliversRider(c *gin.Context) {

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

	result, err := restCtrl.Repo.FetchCompletedDeliversOfRider()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

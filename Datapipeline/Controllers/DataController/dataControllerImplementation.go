package dataController

import (
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

func (data *Controller) SaveConfiguration(ctx *gin.Context) {

	var Config model.Configuration
	if err := ctx.ShouldBindJSON(&Config); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, ctx, "message", "error while binding", "error", err)
		return
	}

	err := data.DriveClient.CreateConnection(Config)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, ctx, "message", "error while connection with the client", "error", err)
		return
	}

	utils.GenerateResponse(http.StatusOK, ctx, "message", "Configuration have been saved in database", "", nil)
	return
}

package dataController

import (
	"fmt"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

func (data *Controller) SourceConfiguration(ctx *gin.Context) {

	var Config model.CombinedStorageConfig
	if err := ctx.ShouldBindJSON(&Config); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, ctx, "message", "error while binding", "error", err)
		return
	}

	err := data.DriveClient.CreateConnection(Config)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, ctx, "message", "error while connection with the client", "", nil)
		return
	}

	err = data.Repo.InsertSourceConfiguration(&Config)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, ctx, "message", "error while connection with the client", "error", err.Error())
		return
	}

	mess := fmt.Sprintf("Source Configuration has been saved in the database and your source config ID is %v", Config.Source.SourcesID)
	utils.GenerateResponse(http.StatusOK, ctx, "message", mess, "", nil)
}

func (data *Controller) DestinationConfiguration(ctx *gin.Context) {

	var Config model.CombinedStorageConfig
	if err := ctx.ShouldBindJSON(&Config); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, ctx, "message", "error while binding", "error", err)
		return
	}

	err := data.DriveClient.CreateConnection(Config)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, ctx, "message", "error while connection with the client", "", nil)
		return
	}

	err = data.Repo.InsertDestinationConfiguration(&Config)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, ctx, "message", "error while connection with the client", "error", err.Error())
		return
	}

	mess := fmt.Sprintf("Destination Configuration has been saved in the database and your destination config ID is %v", Config.Destination.DestinationsID)
	utils.GenerateResponse(http.StatusOK, ctx, "message", mess, "", nil)
}

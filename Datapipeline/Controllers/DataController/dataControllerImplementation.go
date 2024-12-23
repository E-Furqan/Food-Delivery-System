package dataController

import (
	"fmt"
	"log"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

func (data *Controller) SourceConfiguration(ctx *gin.Context) {

	var CombinedStorageConfig model.CombinedStorageConfig
	if err := ctx.ShouldBindJSON(&CombinedStorageConfig); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, ctx, "message", "error while binding", "error", err)
		return
	}

	err := data.DriveClient.CreateConnection(CombinedStorageConfig)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, ctx, "message", "error while connection with the client", "", nil)
		return
	}

	source := utils.CreateSourceObj(CombinedStorageConfig)
	config := utils.CreateConfigObj(CombinedStorageConfig)

	err = data.Repo.InsertSourceConfiguration(&source, &config)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, ctx, "message", "error while connection with the client", "error", err.Error())
		return
	}

	mess := fmt.Sprintf("Source Configuration has been saved in the database and your source config ID is %v", config.SourcesID)
	utils.GenerateResponse(http.StatusOK, ctx, "message", mess, "", nil)
}

func (data *Controller) DestinationConfiguration(ctx *gin.Context) {

	var CombinedStorageConfig model.CombinedStorageConfig
	if err := ctx.ShouldBindJSON(&CombinedStorageConfig); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, ctx, "message", "error while binding", "error", err)
		return
	}

	err := data.DriveClient.CreateConnection(CombinedStorageConfig)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, ctx, "message", "error while connection with the client", "", nil)
		return
	}

	destination := utils.CreateDestinationObj(CombinedStorageConfig)
	config := utils.CreateConfigObj(CombinedStorageConfig)

	err = data.Repo.InsertDestinationConfiguration(&destination, &config)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, ctx, "message", "error while connection with the client", "error", err.Error())
		return
	}

	mess := fmt.Sprintf("Destination Configuration has been saved in the database and your destination config ID is %v", config.DestinationsID)
	utils.GenerateResponse(http.StatusOK, ctx, "message", mess, "", nil)
}

func (data *Controller) CreatePipeline(ctx *gin.Context) {

	var pipeline model.Pipeline
	if err := ctx.ShouldBindJSON(&pipeline); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, ctx, "message", "error while binding", "error", err)
		return
	}

	err := data.Repo.CreatePipeline(pipeline)
	if err != nil {
		utils.GenerateResponse(http.StatusInternalServerError, ctx, "message", "error while creating pipline", "error", err.Error())
		return
	}

	mess := fmt.Sprintf("Your pipeline has been created and your pipeline ID is %v", pipeline.PipelinesID)
	utils.GenerateResponse(http.StatusOK, ctx, "message", mess, "", nil)
}

func (data *Controller) DatapipelineSync(ctx *gin.Context) {

	var pipelineID model.Pipeline
	if err := ctx.ShouldBindJSON(&pipelineID); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, ctx, "message", "error while binding", "error", err)
		return
	}

	pipelineDetails, err := data.Repo.FetchPipelineDetails(pipelineID.PipelinesID)
	if err != nil {
		log.Print("error while fetching error: ", err)
		return
	}

	err = data.WorkFlow.DatapipelineSync(pipelineDetails)
	if err != nil {
		log.Print("error while starting data syn workflow: ", err)
		return
	}

	utils.GenerateResponse(http.StatusOK, ctx, "message", "Data Sync has started", "", nil)
}

package dataController

import (
	"fmt"
	"log"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

func (data *Controller) CreateSourceConfiguration(ctx *gin.Context) {

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

func (data *Controller) CreateDestinationConfiguration(ctx *gin.Context) {

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

func (data *Controller) StartDatapipelineSync(ctx *gin.Context) {

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

func (data *Controller) FetchSourceConfiguration(ctx *gin.Context) {

	var sourceID model.Source
	if err := ctx.ShouldBindJSON(&sourceID); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, ctx, "message", "error while binding", "error", err)
		return
	}

	configDetails, err := data.Repo.FetchConfigSourceDetails(sourceID.SourcesID)
	if err != nil {
		log.Print("error while fetching error: ", err)
		return
	}

	ctx.JSON(http.StatusOK, configDetails)
}

func (data *Controller) FetchDestinationConfiguration(ctx *gin.Context) {

	var destinationID model.Destination
	if err := ctx.ShouldBindJSON(&destinationID); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, ctx, "message", "error while binding", "error", err)
		return
	}

	configDetails, err := data.Repo.FetchConfigDestinationDetails(destinationID.DestinationsID)
	if err != nil {
		log.Print("error while fetching error: ", err)
		return
	}

	ctx.JSON(http.StatusOK, configDetails)
}

func (data *Controller) AddLogs(ctx *gin.Context) {

	var logs model.Log
	if err := ctx.ShouldBindJSON(&logs); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, ctx, "message", "error while binding", "error", err)
		return
	}

	err := data.Repo.AddLogs(logs)
	if err != nil {
		log.Print("error while fetching error: ", err)
		return
	}

	ctx.JSON(http.StatusOK, "Logs have been added")
}

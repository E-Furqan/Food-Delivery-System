package dataController

import (
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

func (data *Controller) SaveConfiguration(ctx *gin.Context) {

	var Config model.Configs
	if err := ctx.ShouldBindJSON(&Config); err != nil {
		utils.GenerateResponse(http.StatusBadRequest, ctx, "message", "error while binding", "error", err)
		return
	}

	// token, err := data.DriveClient.CreateConnection(Config, ctx)
	// if err != nil {
	// 	utils.GenerateResponse(http.StatusInternalServerError, ctx, "message", "error while connection with the client", "", nil)
	// 	return
	// }

	// err = data.Repo.InsertConfiguration(&Config)
	// if err != nil {
	// 	utils.GenerateResponse(http.StatusInternalServerError, ctx, "message", "error while connection with the client", "error", err.Error())
	// 	return
	// }

	// mess := fmt.Sprintf("Configuration has been saved in the database and your config ID is %v", Config.ConfigID)
	// utils.GenerateResponse(http.StatusOK, ctx, "message", mess, "", nil)
}

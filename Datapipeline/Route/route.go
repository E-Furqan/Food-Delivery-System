package route

import (
	driveClient "github.com/E-Furqan/Food-Delivery-System/Client/DriveClient"
	dataController "github.com/E-Furqan/Food-Delivery-System/Controllers/DataController"
	"github.com/gin-gonic/gin"
)

func User_routes(DataCon dataController.DataControllerInterface, DriveClient driveClient.DriveClientInterface, server *gin.Engine) {

	pipeline := server.Group("/pipeline")
	pipeline.POST("/source/configuration", DataCon.SaveConfiguration)
	pipeline.GET("/callback", DriveClient.HandleOAuth2Callback)

}

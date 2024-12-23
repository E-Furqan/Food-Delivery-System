package route

import (
	driveClient "github.com/E-Furqan/Food-Delivery-System/Client/DriveClient"
	dataController "github.com/E-Furqan/Food-Delivery-System/Controllers/DataController"
	"github.com/gin-gonic/gin"
)

func User_routes(DataCon dataController.DataControllerInterface, DriveClient driveClient.DriveClientInterface, server *gin.Engine) {

	pipeline := server.Group("/pipeline")
	pipeline.POST("/source/configuration", DataCon.CreateSourceConfiguration)
	pipeline.POST("/destination/configuration", DataCon.CreateDestinationConfiguration)
	pipeline.POST("/create/pipeline", DataCon.CreatePipeline)
	pipeline.POST("/data/sync", DataCon.StartDatapipelineSync)
	pipeline.GET("/fetch/source/configuration", DataCon.FetchSourceConfiguration)
	pipeline.GET("/fetch/destination/configuration", DataCon.FetchDestinationConfiguration)

}

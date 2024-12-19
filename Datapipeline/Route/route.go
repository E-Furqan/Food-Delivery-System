package route

import (
	dataController "github.com/E-Furqan/Food-Delivery-System/Controllers/DataController"
	"github.com/gin-gonic/gin"
)

func User_routes(DataCon dataController.DataControllerInterface, server *gin.Engine) {

	pipeline := server.Group("/pipeline")
	pipeline.POST("/source/configuration", DataCon.SaveConfiguration)

}

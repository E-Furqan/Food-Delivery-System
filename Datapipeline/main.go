package main

import (
	driveClient "github.com/E-Furqan/Food-Delivery-System/Client/DriveClient"
	dataController "github.com/E-Furqan/Food-Delivery-System/Controllers/DataController"
	config "github.com/E-Furqan/Food-Delivery-System/DatabaseConfig"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	route "github.com/E-Furqan/Food-Delivery-System/Route"
	"github.com/gin-gonic/gin"
)

func main() {
	DatabaseEnv := environmentVariable.ReadDatabaseEnv()

	databaseConfig := config.NewDatabase(DatabaseEnv)
	db := databaseConfig.Connection()

	var repo database.RepositoryInterface = database.NewRepository(db)

	var DriveClient driveClient.DriveClientInterface = driveClient.NewClient()
	var DataController dataController.DataControllerInterface = dataController.NewController(repo, DriveClient)

	server := gin.Default()

	route.User_routes(DataController, DriveClient, server)

	server.Run(":8085")

}

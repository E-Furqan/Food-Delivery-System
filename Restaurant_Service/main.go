package main

import (
	Authenticator "github.com/E-Furqan/Food-Delivery-System/Authentication"
	ClientPackage "github.com/E-Furqan/Food-Delivery-System/Client"
	"github.com/E-Furqan/Food-Delivery-System/Controllers/ItemController"
	"github.com/E-Furqan/Food-Delivery-System/Controllers/RestaurantController"
	config "github.com/E-Furqan/Food-Delivery-System/DatabaseConfig"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	route "github.com/E-Furqan/Food-Delivery-System/Routes"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

func main() {
	envVar := environmentVariable.ReadEnv()
	config.SetEnvValue(envVar)
	db := config.Connection()
	utils.SetEnvValue(envVar)
	Authenticator.SetEnvValue(envVar)

	client := ClientPackage.NewClient()
	client.SetEnvValue(envVar)

	repo := database.NewRepository(db)
	ctrl := RestaurantController.NewController(repo, client)
	ItemController := ItemController.NewController(repo)

	server := gin.Default()
	route.User_routes(ctrl, ItemController, server)
	server.Run(":8082")
}

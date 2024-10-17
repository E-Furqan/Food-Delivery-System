package main

import (
	authenticator "github.com/E-Furqan/Food-Delivery-System/Authentication"
	controller "github.com/E-Furqan/Food-Delivery-System/Controllers/RestaurantController"
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
	authenticator.SetEnvValue(envVar)

	repo := database.NewRepository(db)
	ctrl := controller.NewController(repo)

	server := gin.Default()
	route.User_routes(ctrl, server)
	server.Run(":8084")
}

package main

import (
	ClientPackage "github.com/E-Furqan/Food-Delivery-System/Client"
	"github.com/E-Furqan/Food-Delivery-System/Controllers/ItemController"
	"github.com/E-Furqan/Food-Delivery-System/Controllers/RestaurantController"
	config "github.com/E-Furqan/Food-Delivery-System/DatabaseConfig"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	route "github.com/E-Furqan/Food-Delivery-System/Routes"
	"github.com/gin-gonic/gin"
)

func main() {
	envVar := environmentVariable.ReadEnv()
	config.SetEnvValue(envVar)
	db := config.Connection()
	Middleware.SetEnvValue(envVar)

	client := ClientPackage.NewClient()
	client.SetEnvValue(envVar)

	repo := database.NewRepository(db)
	ctrl := RestaurantController.NewController(repo, client)
	ItemController := ItemController.NewController(repo)
	middle := Middleware.NewMiddleware(client)

	server := gin.Default()
	route.Restaurant_routes(ctrl, ItemController, middle, server)
	server.Run(":8082")
}

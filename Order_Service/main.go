package main

import (
	ClientPackage "github.com/E-Furqan/Food-Delivery-System/Client"
	OrderControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/OrderController"
	config "github.com/E-Furqan/Food-Delivery-System/DatabaseConfig"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/E-Furqan/Food-Delivery-System/Routes"
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
	OrderController := OrderControllers.NewController(repo, client)
	server := gin.Default()
	Routes.Order_routes(OrderController, server)
	server.Run(":8081")
}

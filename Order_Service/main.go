package main

import (
	OrderControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/OrderController"
	config "github.com/E-Furqan/Food-Delivery-System/DatabaseConfig"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/E-Furqan/Food-Delivery-System/Routes"
	"github.com/gin-gonic/gin"
)

func main() {
	envVar := environmentVariable.ReadEnv()
	config.SetEnvValue(envVar)
	db := config.Connection()

	repo := database.NewRepository(db)
	OrderController := OrderControllers.NewController(repo)

	server := gin.Default()
	Routes.User_routes(OrderController, server)
	server.Run(":8084")
}

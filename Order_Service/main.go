package main

import (
	RestaurantClient "github.com/E-Furqan/Food-Delivery-System/Client"
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
	config := config.NewDatabase(envVar)
	db := config.Connection()
	middle := Middleware.NewMiddleware(&envVar)
	ResClient := RestaurantClient.NewClient(&envVar)

	repo := database.NewRepository(db)
	OrderController := OrderControllers.NewController(repo, ResClient)
	server := gin.Default()
	Routes.Order_routes(OrderController, middle, server)
	server.Run(":8081")
}

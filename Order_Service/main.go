package main

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/RestaurantClient"
	OrderControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/OrderController"
	config "github.com/E-Furqan/Food-Delivery-System/DatabaseConfig"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/E-Furqan/Food-Delivery-System/Routes"
	"github.com/gin-gonic/gin"
)

func main() {

	DatabaseConfigEnv := environmentVariable.ReadDatabaseConfigEnv()
	RestaurantClientEnv := environmentVariable.ReadRestaurantClientEnv()
	MiddlewareEnv := environmentVariable.ReadMiddlewareEnv()

	config := config.NewDatabase(DatabaseConfigEnv)
	db := config.Connection()

	var repo database.RepositoryInterface = database.NewRepository(db)

	middle := Middleware.NewMiddleware(&MiddlewareEnv)

	var ResClient RestaurantClient.RestaurantClientInterface = RestaurantClient.NewClient(&RestaurantClientEnv)
	var OrderCtrl OrderControllers.OrderControllerInterface = OrderControllers.NewController(repo, ResClient)

	server := gin.Default()

	Routes.Order_routes(OrderCtrl, middle, server)

	server.Run(":8081")
}

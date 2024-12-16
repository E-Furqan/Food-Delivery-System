package main

import (
	activity "github.com/E-Furqan/Food-Delivery-System/Activity"
	"github.com/E-Furqan/Food-Delivery-System/Client/EmailClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/RestaurantClient"
	"github.com/E-Furqan/Food-Delivery-System/Controllers/orderControllers"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	worker "github.com/E-Furqan/Food-Delivery-System/Worker"
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
	Routes "github.com/E-Furqan/Food-Delivery-System/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// DatabaseEnv := environmentVariable.ReadDatabaseEnv()
	OrderClientEnv := environmentVariable.ReadOrderClientEnv()
	RestaurantClientEnv := environmentVariable.ReadRestaurantClientEnv()

	// databaseConfig := config.NewDatabase(DatabaseEnv)
	// db := databaseConfig.Connection()

	// var repo database.RepositoryInterface = database.NewRepository(db)

	var OrdClient OrderClient.OrdClientInterface = OrderClient.NewClient(OrderClientEnv)
	var emailClient EmailClient.EmailClientInterface = EmailClient.NewClient()
	var restaurantClient RestaurantClient.RestaurantClientInterface = RestaurantClient.NewClient(RestaurantClientEnv)

	var activity_var activity.ActivityInterface = activity.NewController(OrdClient, emailClient, restaurantClient)
	var workFlow workflows.WorkflowInterface = workflows.NewController(activity_var)
	var worker_var worker.WorkerInterface = worker.NewController(activity_var, workFlow)
	var orderController orderControllers.OrderControllerInterface = orderControllers.NewController(workFlow, activity_var)

	server := gin.Default()
	Routes.Order_routes(orderController, server)
	// Start the worker in a goroutine
	go func() {
		worker_var.WorkerUserStart()
	}()
	server.Run(":8084")

}

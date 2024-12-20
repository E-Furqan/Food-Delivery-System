package main

import (
	activity "github.com/E-Furqan/Food-Delivery-System/Activity"
	"github.com/E-Furqan/Food-Delivery-System/Client/EmailClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/RestaurantClient"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	worker "github.com/E-Furqan/Food-Delivery-System/Worker"
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
	Routes "github.com/E-Furqan/Food-Delivery-System/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	OrderClientEnv := environmentVariable.ReadOrderClientEnv()
	RestaurantClientEnv := environmentVariable.ReadRestaurantClientEnv()
	EmailClientEnv := environmentVariable.ReadEmailClientEnv()

	var OrdClient OrderClient.OrdClientInterface = OrderClient.NewClient(OrderClientEnv)
	var emailClient EmailClient.EmailClientInterface = EmailClient.NewClient(EmailClientEnv)
	var restaurantClient RestaurantClient.RestaurantClientInterface = RestaurantClient.NewClient(RestaurantClientEnv)

	var activity_var activity.ActivityInterface = activity.NewController(OrdClient, emailClient, restaurantClient)
	var workFlow workflows.WorkflowInterface = workflows.NewController(activity_var)
	var worker_var worker.WorkerInterface = worker.NewController(activity_var, workFlow)

	server := gin.Default()
	Routes.Workflow_routes(workFlow, server)

	go func() {
		worker_var.WorkerUserStart()
	}()
	server.Run(":8088")

}

package main

import (
	activity "github.com/E-Furqan/Food-Delivery-System/Activity"
	"github.com/E-Furqan/Food-Delivery-System/Client/EmailClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/RestaurantClient"
	userClient "github.com/E-Furqan/Food-Delivery-System/Client/UserClient"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	worker "github.com/E-Furqan/Food-Delivery-System/Worker"
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
	"github.com/E-Furqan/Food-Delivery-System/controllers"
	Routes "github.com/E-Furqan/Food-Delivery-System/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	OrderClientEnv := environmentVariable.ReadOrderClientEnv()
	RestaurantClientEnv := environmentVariable.ReadRestaurantClientEnv()
	EmailClientEnv := environmentVariable.ReadEmailClientEnv()
	UserClientEnv := environmentVariable.ReadUserClientEnv()

	var OrdClient OrderClient.OrdClientInterface = OrderClient.NewClient(OrderClientEnv)
	var emailClient EmailClient.EmailClientInterface = EmailClient.NewClient(EmailClientEnv)
	var restaurantClient RestaurantClient.RestaurantClientInterface = RestaurantClient.NewClient(RestaurantClientEnv)
	var UserClient userClient.UserClientInterface = userClient.NewClient(UserClientEnv)

	var activity_var activity.ActivityInterface = activity.NewController(OrdClient, emailClient, restaurantClient, UserClient)
	var workFlow workflows.WorkflowInterface = workflows.NewController(activity_var)
	var worker_var worker.WorkerInterface = worker.NewController(activity_var, workFlow)
	var controller controllers.ControllerInterface = controllers.NewController(workFlow)

	server := gin.Default()
	Routes.Workflow_routes(controller, server)

	go func() {
		worker_var.WorkerUserStart()
	}()
	server.Run(":8088")

}

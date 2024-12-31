package main

import (
	activityPac "github.com/E-Furqan/Food-Delivery-System/Activity"
	datapipelineClient "github.com/E-Furqan/Food-Delivery-System/Client/DatapipelineClient"
	driveClient "github.com/E-Furqan/Food-Delivery-System/Client/DriveClient"
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
	datapipelineClientEnv := environmentVariable.ReadPipelineEnv()

	var OrdClient OrderClient.OrdClientInterface = OrderClient.NewClient(OrderClientEnv)
	var emailClient EmailClient.EmailClientInterface = EmailClient.NewClient(EmailClientEnv)
	var restaurantClient RestaurantClient.RestaurantClientInterface = RestaurantClient.NewClient(RestaurantClientEnv)
	var UserClient userClient.UserClientInterface = userClient.NewClient(UserClientEnv)
	var DatapipelineClient datapipelineClient.DatapipelineClientInterface = datapipelineClient.NewClient(datapipelineClientEnv)
	var DriveClient driveClient.DriveClientInterface = driveClient.NewClient()

	var activity_var activityPac.ActivityInterface = activityPac.NewController(OrdClient, emailClient, restaurantClient, UserClient, DatapipelineClient, DriveClient)
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

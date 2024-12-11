package main

import (
	activity "github.com/E-Furqan/Food-Delivery-System/Activity"
	"github.com/E-Furqan/Food-Delivery-System/Client/AuthClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	RoleController "github.com/E-Furqan/Food-Delivery-System/Controllers/RoleControler"
	UserControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/UserController"
	config "github.com/E-Furqan/Food-Delivery-System/DatabaseConfig"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	route "github.com/E-Furqan/Food-Delivery-System/Route"
	worker "github.com/E-Furqan/Food-Delivery-System/Worker"
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
	"github.com/gin-gonic/gin"
)

func main() {
	DatabaseEnv := environmentVariable.ReadDatabaseEnv()
	OrderClientEnv := environmentVariable.ReadOrderClientEnv()
	AuthClientEnv := environmentVariable.ReadAuthClientEnv()
	MiddlewareEnv := environmentVariable.ReadMiddlewareEnv()

	databaseConfig := config.NewDatabase(DatabaseEnv)
	db := databaseConfig.Connection()

	var repo database.RepositoryInterface = database.NewRepository(db)

	var OrdClient OrderClient.OrdClientInterface = OrderClient.NewClient(OrderClientEnv)
	var AuthClient AuthClient.AuthClientInterface = AuthClient.NewClient(AuthClientEnv)
	var activity_var activity.ActivityInterface = activity.NewController(repo)
	var workFlow workflows.WorkflowInterface = workflows.NewController(repo, activity_var)
	var worker_var worker.WorkerInterface = worker.NewController(repo, activity_var, workFlow)
	// Start the worker in a goroutine
	go func() {
		worker_var.WorkerUserStart()
	}()

	var uCtrl UserControllers.UserControllerInterface
	var rCtrl RoleController.RoleControllerInterface

	uCtrl = UserControllers.NewController(repo, OrdClient, AuthClient, workFlow)
	rCtrl = RoleController.NewController(repo, AuthClient)

	var middle Middleware.MiddlewareInterface = Middleware.NewMiddleware(AuthClient, &MiddlewareEnv)

	server := gin.Default()

	rCtrl.AddDefaultRoles(nil)

	route.User_routes(uCtrl, rCtrl, middle, server)

	server.Run(":8083")

}

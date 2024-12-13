package workflow

import (
	activity "github.com/E-Furqan/Food-Delivery-System/Activity"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	config "github.com/E-Furqan/Food-Delivery-System/DatabaseConfig"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	worker "github.com/E-Furqan/Food-Delivery-System/Worker"
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
	"github.com/gin-gonic/gin"
)

func main() {
	DatabaseEnv := environmentVariable.ReadDatabaseEnv()
	OrderClientEnv := environmentVariable.ReadOrderClientEnv()

	databaseConfig := config.NewDatabase(DatabaseEnv)
	db := databaseConfig.Connection()

	var repo database.RepositoryInterface = database.NewRepository(db)

	var OrdClient OrderClient.OrdClientInterface = OrderClient.NewClient(OrderClientEnv)
	var activity_var activity.ActivityInterface = activity.NewController(repo, OrdClient)
	var workFlow workflows.WorkflowInterface = workflows.NewController(repo, activity_var)
	var worker_var worker.WorkerInterface = worker.NewController(repo, activity_var, workFlow)
	// Start the worker in a goroutine
	go func() {
		worker_var.WorkerUserStart()
	}()

	server := gin.Default()

	server.Run(":8084")

}

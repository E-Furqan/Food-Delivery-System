package main

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/AuthClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	RoleController "github.com/E-Furqan/Food-Delivery-System/Controllers/RoleControler"
	UserControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/UserController"
	config "github.com/E-Furqan/Food-Delivery-System/DatabaseConfig"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	route "github.com/E-Furqan/Food-Delivery-System/Route"
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

	var uCtrl UserControllers.UserControllerInterface
	var rCtrl RoleController.RoleControllerInterface

	uCtrl = UserControllers.NewController(repo, OrdClient, AuthClient)
	rCtrl = RoleController.NewController(repo, AuthClient)

	var middle Middleware.MiddlewareInterface = Middleware.NewMiddleware(AuthClient, &MiddlewareEnv)

	server := gin.Default()

	rCtrl.AddDefaultRoles(nil)

	route.User_routes(uCtrl, rCtrl, middle, server)

	server.Run(":8083")

}

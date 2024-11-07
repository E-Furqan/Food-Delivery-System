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
	repo := database.NewRepository(db)

	OrdClient := OrderClient.NewClient(OrderClientEnv)
	AuthClient := AuthClient.NewClient(AuthClientEnv)

	ctrl := UserControllers.NewController(repo, OrdClient, AuthClient)
	rCtrl := RoleController.NewController(repo, AuthClient)

	middle := Middleware.NewMiddleware(AuthClient, &MiddlewareEnv)

	server := gin.Default()

	rCtrl.AddDefaultRoles(nil)

	route.User_routes(ctrl, rCtrl, middle, server)

	server.Run(":8083")
}

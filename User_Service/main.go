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
	envVar := environmentVariable.ReadEnv()
	databaseConfig := config.NewDatabase(envVar)

	db := databaseConfig.Connection()

	repo := database.NewRepository(db)
	OrdClient := OrderClient.NewClient(envVar)
	AuthClient := AuthClient.NewClient(envVar)

	ctrl := UserControllers.NewController(repo, OrdClient, AuthClient)
	rCtrl := RoleController.NewController(repo)
	middle := Middleware.NewMiddleware(AuthClient, &envVar)

	server := gin.Default()
	rCtrl.AddDefaultRoles(&gin.Context{})

	route.User_routes(ctrl, rCtrl, middle, server)
	server.Run(":8083")
}

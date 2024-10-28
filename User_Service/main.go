package main

import (
	ClientPackage "github.com/E-Furqan/Food-Delivery-System/Client"
	RoleController "github.com/E-Furqan/Food-Delivery-System/Controllers/RoleControler"
	UserControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/UserController"
	config "github.com/E-Furqan/Food-Delivery-System/DatabaseConfig"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	route "github.com/E-Furqan/Food-Delivery-System/Route"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
)

func main() {
	envVar := environmentVariable.ReadEnv()

	config.SetEnvValue(envVar)
	db := config.Connection()
	utils.SetEnvValue(envVar)
	Middleware.SetEnvValue(envVar)

	repo := database.NewRepository(db)
	client := ClientPackage.NewClient()
	client.SetEnvValue(envVar)
	ctrl := UserControllers.NewController(repo, client)
	rCtrl := RoleController.NewController(repo)
	middle := Middleware.NewMiddleware(client)

	server := gin.Default()
	rCtrl.AddDefaultRoles(&gin.Context{})

	route.User_routes(ctrl, rCtrl, middle, server)
	server.Run(":8083")
}

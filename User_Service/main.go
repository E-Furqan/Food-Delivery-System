package main

import (
	authenticator "github.com/E-Furqan/Food-Delivery-System/Authentication"
	ClientPackage "github.com/E-Furqan/Food-Delivery-System/Client"
	RoleController "github.com/E-Furqan/Food-Delivery-System/Controllers/RoleControler"
	UserControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/UserController"
	config "github.com/E-Furqan/Food-Delivery-System/DatabaseConfig"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
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
	authenticator.SetEnvValue(envVar)

	repo := database.NewRepository(db)
	client := ClientPackage.NewClient()
	ctrl := UserControllers.NewController(repo, client)
	rCtrl := RoleController.NewController(repo)

	server := gin.Default()

	rCtrl.AddDefaultRoles(&gin.Context{})

	route.User_routes(ctrl, rCtrl, server)
	server.Run(":8083")
}

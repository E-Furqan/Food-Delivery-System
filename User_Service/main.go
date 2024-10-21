package main

import (
	authenticator "github.com/E-Furqan/Food-Delivery-System/Authentication"
	"github.com/E-Furqan/Food-Delivery-System/Controllers/OrderController"
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
	db := config.Connection()
	envVar := environmentVariable.ReadEnv()
	utils.SetEnvValue(envVar)
	authenticator.SetEnvValue(envVar)

	repo := database.NewRepository(db)
	ctrl := UserControllers.NewController(repo)
	rCtrl := RoleController.NewController(repo)
	orderCtrl := OrderController.NewController(repo)

	server := gin.Default()

	rCtrl.AddDefaultRoles(&gin.Context{})

	route.User_routes(ctrl, rCtrl, orderCtrl, server)
	server.Run(":8084")
}

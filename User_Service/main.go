package main

import (
	authenticator "github.com/E-Furqan/Food-Delivery-System/Authentication"
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
	rctrl := RoleController.NewController(repo)

	server := gin.Default()

	// Invoke AddDefaultRoles function to initialize roles
	rctrl.AddDefaultRoles(&gin.Context{}) // Passing an empty context, can be modified as per requirement

	route.User_routes(ctrl, rctrl, server)
	server.Run(":8083")
}

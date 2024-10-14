package main

import (
	RoleController "github.com/E-Furqan/Food-Delivery-System/Controllers/RoleControler"
	UserControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/UserController"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	config "github.com/E-Furqan/Food-Delivery-System/database_config"
	"github.com/E-Furqan/Food-Delivery-System/route"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.Connection()

	repo := database.NewRepository(db)
	ctrl := UserControllers.NewController(repo)
	rctrl := RoleController.NewController(repo)

	server := gin.Default()

	// Invoke AddDefaultRoles function to initialize roles
	rctrl.AddDefaultRoles(&gin.Context{}) // Passing an empty context, can be modified as per requirement

	route.User_routes(ctrl, rctrl, server)
	server.Run(":8083")
}

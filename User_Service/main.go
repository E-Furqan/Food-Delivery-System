package main

import (
	UserControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/UserController"
	RoleController "github.com/E-Furqan/Food-Delivery-System/Controllers/UserController/RoleControler"
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
	route.User_routes(ctrl, rctrl, server)
	server.Run(":8082")
}

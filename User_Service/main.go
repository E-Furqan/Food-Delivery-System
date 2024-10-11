package main

import (
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
	rctrl ;= 
	server := gin.Default()
	route.User_routes(ctrl, server)
	server.Run(":8082")
}

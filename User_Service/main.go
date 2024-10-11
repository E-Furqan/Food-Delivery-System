package main

import (
	controllers "github.com/E-Furqan/Food-Delivery-System/Interfaces/Controllers"
	database "github.com/E-Furqan/Food-Delivery-System/Interfaces/Repositories"
	config "github.com/E-Furqan/Food-Delivery-System/database_config"
	"github.com/E-Furqan/Food-Delivery-System/route"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.Connection()

	repo := database.NewRepository(db)
	// Initialize the controller with the repository
	ctrl := controllers.NewController(repo)
	server := gin.Default()
	route.User_routes(ctrl, server)
	server.Run(":8082")
}

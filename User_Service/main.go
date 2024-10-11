package main

import (
	config "github.com/E-Furqan/Food-Delivery-System/database_config"
	controllers "github.com/E-Furqan/Food-Delivery-System/handelers/Controllers"
	database "github.com/E-Furqan/Food-Delivery-System/handelers/Repositories"
	"github.com/E-Furqan/Food-Delivery-System/route"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.Connection()

	repo := database.NewRepository(db)
	ctrl := controllers.NewController(repo)
	server := gin.Default()
	route.User_routes(ctrl, server)
	server.Run(":8082")
}

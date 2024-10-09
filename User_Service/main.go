package main

import (
	controllers "github.com/E-Furqan/Food-Delivery-System/Controllers"
	config "github.com/E-Furqan/Food-Delivery-System/database_config"
	"github.com/E-Furqan/Food-Delivery-System/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Connection()
	server := gin.Default()

	server.POST("/Register", controllers.Register)
	server.GET("/Getuser", controllers.Getuser)
	server.GET("/Getrole", controllers.Getrole)
	server.POST("/login", controllers.Login)

	protected := server.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.PATCH("/updaterole", controllers.Change_Role)
	}
	server.Run(":8081")
}

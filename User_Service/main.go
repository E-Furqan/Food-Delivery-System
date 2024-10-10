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
	server.POST("/Login", controllers.Login)

	protected := server.Group("/protected")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.PATCH("/Update_role", controllers.Update_Role)
		protected.PATCH("/Update_user", controllers.Update_user)
		protected.DELETE("/Delete_user", controllers.Delete_user)
		protected.DELETE("/Delete_role", controllers.Delete_role)
	}
	server.Run(":8081")
}

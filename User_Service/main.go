package main

import (
	controllers "github.com/E-Furqan/Food-Delivery-System/Controllers"
	config "github.com/E-Furqan/Food-Delivery-System/database_config"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Connection()
	server := gin.Default()
	server.POST("/login", controllers.Login)
	server.POST("/Register", controllers.Register)
	server.GET("/Getuser", controllers.Getuser)
	server.PATCH("/updaterole", controllers.Change_Role)
}

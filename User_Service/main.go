package main

import (
	controllers "github.com/E-Furqan/Food-Delivery-System/Interfaces/Controllers"
	database "github.com/E-Furqan/Food-Delivery-System/Interfaces/Repositories"
	config "github.com/E-Furqan/Food-Delivery-System/database_config"
	"github.com/E-Furqan/Food-Delivery-System/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.Connection()

	repo := database.NewRepository(db)
	// Initialize the controller with the repository
	ctrl := controllers.NewController(repo)

	server := gin.Default()

	server.POST("/app/Register", ctrl.Register)
	server.GET("/app/Getuser", ctrl.Get_user)
	server.GET("/app/Getrole", ctrl.Get_role)
	server.POST("/app/Login", ctrl.Login)

	protected := server.Group("/app")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.PATCH("/role/update", ctrl.Update_Role)
		protected.PATCH("/user/update", ctrl.Update_user)
		protected.DELETE("/user/delete", ctrl.Delete_user)
		protected.DELETE("/role/delete", ctrl.Delete_role)
	}
	server.Run(":8081")
}

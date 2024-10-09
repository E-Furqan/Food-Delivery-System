package controllers

import (
	"net/http"

	data "github.com/E-Furqan/Food-Delivery-System/Data"
	config "github.com/E-Furqan/Food-Delivery-System/database_config"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var reg_data data.User
	if err := c.ShouldBindJSON(&reg_data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(reg_data)
	c.JSON(http.StatusCreated, reg_data)
}

func getuser(c *gin.Context) {
	var user_data []data.User
}

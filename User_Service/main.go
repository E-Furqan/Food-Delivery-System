package main

import (
	config "github.com/E-Furqan/Food-Delivery-System/database_config"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Connection()
	server := gin.Default()

}

package main

import (
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/EnvironmentVariable"
	route "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Route"
	utils "github.com/E-Furqan/Food-Delivery-System/Authentication_Service/Utils"
	"github.com/gin-gonic/gin"
)

func main() {
	envVar := environmentVariable.ReadEnv()

	utils.SetEnvValue(envVar)
	server := gin.Default()
	route.Auth_routes(server)
	server.Run(":8084")
}

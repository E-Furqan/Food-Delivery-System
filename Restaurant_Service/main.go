package main

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/AuthClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	"github.com/E-Furqan/Food-Delivery-System/Controllers/ItemController"
	"github.com/E-Furqan/Food-Delivery-System/Controllers/RestaurantController"
	config "github.com/E-Furqan/Food-Delivery-System/DatabaseConfig"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	route "github.com/E-Furqan/Food-Delivery-System/Routes"
	"github.com/gin-gonic/gin"
)

func main() {
	DatabaseEnv := environmentVariable.ReadDatabaseEnv()
	OrderClientEnv := environmentVariable.ReadOrderClientEnv()
	AuthClientEnv := environmentVariable.ReadAuthClientEnv()
	MiddlewareEnv := environmentVariable.ReadMiddlewareEnv()

	databaseConfig := config.NewDatabase(DatabaseEnv)
	db := databaseConfig.Connection()

	var repo database.RepositoryInterface = database.NewRepository(db)

	var OrdClient OrderClient.OrdClientInterface = OrderClient.NewClient(OrderClientEnv)
	var AuthClient AuthClient.AuthClientInterface = AuthClient.NewClient(AuthClientEnv)

	var RestaurantCtrl RestaurantController.RestaurantControllerInterface
	var ItemCtrl ItemController.ItemControllerInterface
	RestaurantCtrl = RestaurantController.NewController(repo, OrdClient, AuthClient)
	ItemCtrl = ItemController.NewController(repo)

	var middleware Middleware.MiddlewareInterface = Middleware.NewMiddleware(AuthClient, &MiddlewareEnv)

	server := gin.Default()
	route.Restaurant_routes(RestaurantCtrl, ItemCtrl, middleware, server)
	server.Run(":8082")
}

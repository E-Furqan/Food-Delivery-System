package main

import (
	"github.com/E-Furqan/Food-Delivery-System/Client/RestaurantClient"
	CustomerController "github.com/E-Furqan/Food-Delivery-System/Controllers/CustomerContoller"
	RiderController "github.com/E-Furqan/Food-Delivery-System/Controllers/DeliverRiderController"
	OrderControllers "github.com/E-Furqan/Food-Delivery-System/Controllers/OrderController"
	"github.com/E-Furqan/Food-Delivery-System/Controllers/RestaurantController"
	config "github.com/E-Furqan/Food-Delivery-System/DatabaseConfig"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	"github.com/E-Furqan/Food-Delivery-System/Middleware"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/E-Furqan/Food-Delivery-System/Routes"
	"github.com/gin-gonic/gin"
)

func main() {

	DatabaseConfigEnv := environmentVariable.ReadDatabaseConfigEnv()
	RestaurantClientEnv := environmentVariable.ReadRestaurantClientEnv()
	MiddlewareEnv := environmentVariable.ReadMiddlewareEnv()

	config := config.NewDatabase(DatabaseConfigEnv)
	db := config.Connection()

	var repo database.RepositoryInterface = database.NewRepository(db)

	var middle Middleware.MiddlewareInterface = Middleware.NewMiddleware(&MiddlewareEnv)

	var ResClient RestaurantClient.RestaurantClientInterface = RestaurantClient.NewClient(&RestaurantClientEnv)
	var OrderCtrl OrderControllers.OrderControllerInterface = OrderControllers.NewController(repo, ResClient)
	var restCtrl RestaurantController.RestaurantControllerInterface = RestaurantController.NewController(repo)
	var cusCtrl CustomerController.CustomerControllerInterface = CustomerController.NewController(repo)
	var riderCtrl RiderController.RiderControllerInterface = RiderController.NewController(repo)

	server := gin.Default()

	Routes.Order_routes(OrderCtrl, restCtrl, cusCtrl, riderCtrl, middle, server)

	server.Run(":8081")
}

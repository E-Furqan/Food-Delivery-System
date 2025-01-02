package activityPac

import (
	"context"

	datapipelineClient "github.com/E-Furqan/Food-Delivery-System/Client/DatapipelineClient"
	driveClient "github.com/E-Furqan/Food-Delivery-System/Client/DriveClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/EmailClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/OrderClient"
	"github.com/E-Furqan/Food-Delivery-System/Client/RestaurantClient"
	userClient "github.com/E-Furqan/Food-Delivery-System/Client/UserClient"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

type Activity struct {
	OrderClient        OrderClient.OrdClientInterface
	Email              EmailClient.EmailClientInterface
	ResClient          RestaurantClient.RestaurantClientInterface
	UserClient         userClient.UserClientInterface
	DatapipelineClient datapipelineClient.DatapipelineClientInterface
	DriveClient        driveClient.DriveClientInterface
}

func NewController(orderClient OrderClient.OrdClientInterface,
	email EmailClient.EmailClientInterface, resClient RestaurantClient.RestaurantClientInterface,
	userClient userClient.UserClientInterface,
	datapipeline datapipelineClient.DatapipelineClientInterface,
	driveClient driveClient.DriveClientInterface) *Activity {
	return &Activity{
		OrderClient:        orderClient,
		Email:              email,
		ResClient:          resClient,
		UserClient:         userClient,
		DatapipelineClient: datapipeline,
		DriveClient:        driveClient,
	}
}

type ActivityInterface interface {
	GetItems(order model.CombineOrderItem, token string) ([]model.Items, error)
	CalculateBill(CombineOrderItem model.CombineOrderItem, items []model.Items) (float64, error)
	CreateOrder(order model.CombineOrderItem, token string) (model.UpdateOrder, error)
	SendEmail(orderID uint, orderStatus string, token string, userEmail string) (string, error)
	CheckOrderStatus(orderID uint, token string) (string, error)
	FetchUserEmail(token string) (*model.UserEmail, error)
	FetchSourceConfiguration(source model.Source) (model.Config, error)
	FetchDestinationConfiguration(destination model.Destination) (model.Config, error)
	CreateSourceToken(source model.Config) (string, error)
	CreateDestinationToken(destination model.Config) (string, error)
	AddLogs(counter model.FileCounter, PipelinesID int) error
	CountFilesInFolder(sourceToken string, sourceConfig model.Config, folderID string) (int, error)
	MoveDataFromSourceToDestination(ctx context.Context, sourceToken string, destinationToken string,
		sourceFolderUrl string, destinationFolderUrl string, sourceConfig model.Config, batchSize int) (model.FileCounter, error)

	MoveBatchActivity(ctx context.Context, sourceToken string, destinationToken string, sourceConfig model.Config,
		destinationConfig model.Config, sourceFolderID string, destinationFolderID string, startIndex int, endIndex int) error
}

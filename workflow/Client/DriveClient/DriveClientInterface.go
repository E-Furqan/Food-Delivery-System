package driveClient

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"google.golang.org/api/drive/v3"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

type DriveClientInterface interface {
	CreateConnection(config model.Config) (*drive.Service, error)
}

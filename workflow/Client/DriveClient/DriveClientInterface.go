package driveClient

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

type DriveClientInterface interface {
	CreateConnection(config model.Config) error
}

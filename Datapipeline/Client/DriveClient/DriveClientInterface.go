package driveClient

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
)

type Client struct {
	codeChan chan string
	Repo     database.RepositoryInterface
}

func NewClient(repo database.RepositoryInterface) *Client {
	return &Client{
		codeChan: make(chan string),
		Repo:     repo,
	}
}

type DriveClientInterface interface {
	CreateSourceConnection(config model.CombinedSourceStorageConfig) error
	CreateDestinationConnection(config model.CombinedDestinationStorageConfig) error
}

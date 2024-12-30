package driveClient

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
)

type Client struct {
	Repo database.RepositoryInterface
}

func NewClient(repo database.RepositoryInterface) *Client {
	return &Client{
		Repo: repo,
	}
}

type DriveClientInterface interface {
	CreateConnection(config model.Config) error
}

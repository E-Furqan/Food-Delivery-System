package driveClient

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"golang.org/x/oauth2"
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
	loadToken(file string) (*oauth2.Token, error)
	CreateConnection(config model.Config) error
}

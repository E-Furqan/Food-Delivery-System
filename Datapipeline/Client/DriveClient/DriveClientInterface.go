package driveClient

import (
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"golang.org/x/oauth2"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

type DriveControllerInterface interface {
	getClient(config *oauth2.Config) (*http.Client, error)
	loadToken(file string) (*oauth2.Token, error)
	saveToken(path string, token *oauth2.Token) error
	getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error)
	CreateConnection(config model.Configuration) error
}

package driveClient

import (
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/gin-gonic/gin"
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
	getClient(config *oauth2.Config, ctx *gin.Context) (*http.Client, *oauth2.Token, error)
	loadToken(file string) (*oauth2.Token, error)
	getTokenFromWeb(config *oauth2.Config, ctx *gin.Context) (*oauth2.Token, error)
	CreateConnection(config model.Configs, ctx *gin.Context) (*oauth2.Token, error)
	HandleOAuth2Callback(ctx *gin.Context)
}

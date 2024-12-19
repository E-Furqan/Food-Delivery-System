package driveClient

import (
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type Client struct {
	codeChan chan string
}

func NewClient() *Client {
	return &Client{
		codeChan: make(chan string),
	}
}

type DriveClientInterface interface {
	getClient(config *oauth2.Config, ctx *gin.Context) (*http.Client, error)
	loadToken(file string) (*oauth2.Token, error)
	saveToken(path string, token *oauth2.Token) error
	getTokenFromWeb(config *oauth2.Config, ctx *gin.Context) (*oauth2.Token, error)
	CreateConnection(config model.Configuration, ctx *gin.Context) error
	HandleOAuth2Callback(ctx *gin.Context)
}

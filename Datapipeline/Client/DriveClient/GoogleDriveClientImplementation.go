package driveClient

import (
	"encoding/json"
	"os"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func (driveClient *Client) loadToken(filePath string) (*oauth2.Token, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var token oauth2.Token
	err = json.NewDecoder(file).Decode(&token)

	return &token, err
}

func (driveClient *Client) CreateConnection(config model.Configs, ctx *gin.Context) {

}

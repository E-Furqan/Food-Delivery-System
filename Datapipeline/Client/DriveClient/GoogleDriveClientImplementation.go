package driveClient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func (driveClient *Client) getClient(config *oauth2.Config, ctx *gin.Context) (*http.Client, error) {
	tokenFile := "token" + config.ClientID + ".json"

	token, err := driveClient.loadToken(tokenFile)
	if err != nil {
		token, err := driveClient.getTokenFromWeb(config, ctx)
		if err != nil {
			return &http.Client{}, err
		}
		driveClient.saveToken(tokenFile, token)
		return config.Client(ctx, token), nil
	}
	return config.Client(ctx, token), nil
}

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

func (driveClient *Client) saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	file, err := os.Create(path)

	if err != nil {
		log.Printf("Unable to cache oauth token: %v", err)
		return err
	}
	defer file.Close()

	json.NewEncoder(file).Encode(token)

	return nil
}

func (driveClient *Client) getTokenFromWeb(config *oauth2.Config, ctx *gin.Context) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	err := utils.OpenBrowser(authURL)
	if err != nil {
		log.Printf("Unable to open browser: %v", err)
	}
	fmt.Printf("Go to the following link in your browser and get the code:\n%v\n", authURL)

	select {
	case code := <-driveClient.codeChan:

		token, err := config.Exchange(context.Background(), code)

		if err != nil {
			log.Printf("Unable to retrieve token from web: %v", err)
			return &oauth2.Token{}, err
		}

		return token, nil
	case <-time.After(2 * time.Minute):
		return &oauth2.Token{}, fmt.Errorf("timeout waiting for authorization code")
	}
}

func (driveClient *Client) CreateConnection(config model.Configs, ctx *gin.Context) error {

	oauthConfig := utils.CreateAuthObj(config)
	client, err := driveClient.getClient(oauthConfig, ctx)
	log.Print("CLIENt:  ", client)
	if err != nil {
		log.Print("Error while fetching client: ", err)
		return err
	}

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Unable to create Drive client: %v", err)
	}

	files, err := srv.Files.List().PageSize(10).Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Printf("Unable to retrieve files: %v", err)
		return err
	}

	log.Println("Files:", len(files.Files))

	return nil
}

func (driveClient *Client) HandleOAuth2Callback(ctx *gin.Context) {
	code := ctx.DefaultQuery("code", "")
	state := ctx.DefaultQuery("state", "")
	if code == "" || state == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state"})
		return
	}

	driveClient.codeChan <- code

	ctx.JSON(http.StatusOK, gin.H{"message": "Authorization successful! You can return to the application."})
}

package driveClient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
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

func (driveClient *Client) CreateConnection(config model.Config) error {

	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: config.TokenURI,
		},
	}

	token, err := oauthConfig.TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: config.RefreshToken,
	}).Token()
	if err != nil {
		log.Print(err)
		return fmt.Errorf("failed to generate token: %w", err)
	}

	httpClient := oauthConfig.Client(context.Background(), token)

	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return fmt.Errorf("failed to create Google Drive service: %w", err)
	}
	fileList, err := driveService.Files.List().Do()
	if err != nil {
		log.Printf("Unable to retrieve files: %v", err)
		return fmt.Errorf("unable to retrieve files: %w", err)
	}

	log.Print("total files in drive: ", len(fileList.Files))

	return nil
}

package driveClient

import (
	"context"
	"fmt"
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func (driveClient *Client) CreateSourceConnection(config model.CombinedSourceStorageConfig) error {

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

func (driveClient *Client) CreateDestinationConnection(config model.CombinedDestinationStorageConfig) error {

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

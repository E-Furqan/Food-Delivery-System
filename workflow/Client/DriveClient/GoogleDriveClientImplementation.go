package driveClient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func (driveClient *Client) CreateToken(config model.Config) (string, error) {
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
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("failed to serialize token: %w", err)
	}

	return string(tokenJSON), nil
}

func (driveClient *Client) CreateConnection(tokenJSON string, sourceConfig model.Config) (drive.Service, error) {
	var token oauth2.Token
	err := json.Unmarshal([]byte(tokenJSON), &token)
	if err != nil {
		return drive.Service{}, fmt.Errorf("failed to deserialize token: %w", err)
	}

	oauthConfig := &oauth2.Config{
		ClientID:     sourceConfig.ClientID,
		ClientSecret: sourceConfig.ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: sourceConfig.TokenURI,
		},
	}
	httpClient := oauthConfig.Client(context.Background(), &token)

	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return drive.Service{}, fmt.Errorf("failed to create Google Drive service: %w", err)
	}

	return *driveService, nil

}

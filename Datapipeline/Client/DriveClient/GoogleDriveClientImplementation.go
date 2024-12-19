package driveClient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func (driveController *Controller) getClient(config *oauth2.Config) (*http.Client, error) {
	tokenFile := "token.json"

	token, err := driveController.loadToken(tokenFile)
	if err != nil {
		token, err = driveController.getTokenFromWeb(config)
		if err != nil {
			return &http.Client{}, err
		}
		driveController.saveToken(tokenFile, token)
	}
	log.Print("token: ", token)
	return config.Client(context.Background(), token), err
}

func (driveController *Controller) loadToken(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var token oauth2.Token
	err = json.NewDecoder(f).Decode(&token)
	return &token, err
}

func (driveController *Controller) saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Printf("Unable to cache oauth token: %v", err)
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}

func (driveController *Controller) getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser and get the code:\n%v\n", authURL)

	var authCode string
	fmt.Print("Enter the authorization code: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Printf("Unable to read authorization code: %v", err)
		return &oauth2.Token{}, err
	}

	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Printf("Unable to retrieve token from web: %v", err)
		return &oauth2.Token{}, err
	}
	return token, err
}

func (driveController *Controller) CreateConnection(config model.Configuration) error {

	oauthConfig := utils.CreateAuthObj(config)

	client, err := driveController.getClient(oauthConfig)
	if err != nil {
		log.Print("Error while fetching client: ", client)
		return err
	}

	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
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

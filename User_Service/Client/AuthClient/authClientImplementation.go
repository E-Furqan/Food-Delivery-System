package AuthClient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
)

func (client *AuthClient) GenerateToken(input model.UserClaim) (*model.Tokens, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		log.Print("1")
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}

	url, err := utils.CreateUrl(client.AuthClientEnv.BASE_URL,
		client.AuthClientEnv.AUTH_PORT,
		client.AuthClientEnv.GENERATE_TOKEN_URL)

	if err != nil {
		log.Print("2")
		return nil, fmt.Errorf("error: %v", err)
	}

	req, err := utils.CreateRequest(url, jsonData, "POST")
	if err != nil {
		log.Print("3")
		return nil, fmt.Errorf("error: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print("4")
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Print("5")
		return nil, fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	var tokens model.Tokens
	limit := int64(1 << 20)
	if err := json.NewDecoder(io.LimitReader(resp.Body, limit)).Decode(&tokens); err != nil {
		log.Print("6")
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &tokens, nil
}

func (client *AuthClient) RefreshToken(input model.RefreshToken) (*model.Tokens, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}

	url, err := utils.CreateUrl(client.AuthClientEnv.BASE_URL,
		client.AuthClientEnv.AUTH_PORT,
		client.AuthClientEnv.REFRESH_TOKEN_URL)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	req, err := utils.CreateRequest(url, jsonData, "POST")
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	var tokens model.Tokens
	limit := int64(1 << 20)
	if err := json.NewDecoder(io.LimitReader(resp.Body, limit)).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &tokens, nil
}

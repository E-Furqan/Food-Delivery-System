package AuthClient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
)

func (authClient *AuthClient) GenerateToken(input model.RestaurantClaim) (*model.Tokens, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}

	url, err := utils.CreateUrl(authClient.AuthClientEnv.BASE_URL, authClient.AuthClientEnv.AUTH_PORT, authClient.AuthClientEnv.GENERATE_TOKEN_URL)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	req, err := utils.CreateRequest(url, jsonData, "POST")
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	HttpClient := utils.CreateHTTPClient()
	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to generate token: received HTTP %d", resp.StatusCode)
	}

	var tokens model.Tokens
	limit := int64(1 << 20) //limiting the size of response to 1MB
	if err := json.NewDecoder(io.LimitReader(resp.Body, limit)).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &tokens, nil
}

func (authClient *AuthClient) RefreshToken(input model.RefreshToken) (*model.Tokens, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}

	url, err := utils.CreateUrl(authClient.AuthClientEnv.BASE_URL, authClient.AuthClientEnv.AUTH_PORT, authClient.AuthClientEnv.REFRESH_TOKEN_URL)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	req, err := utils.CreateRequest(url, jsonData, "POST")
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	HttpClient := utils.CreateHTTPClient()
	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to refresh token: received HTTP %d", resp.StatusCode)
	}

	var tokens model.Tokens
	limit := int64(1 << 20) //limiting the size of response to 1MB
	if err := json.NewDecoder(io.LimitReader(resp.Body, limit)).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &tokens, nil
}

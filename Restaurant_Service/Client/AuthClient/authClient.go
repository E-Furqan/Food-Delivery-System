package AuthClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

type AuthClient struct {
	environmentVariable.Environment
}

func NewClient(env environmentVariable.Environment) *AuthClient {
	return &AuthClient{
		Environment: env,
	}
}

func (client *AuthClient) GenerateResponse(input model.RestaurantClaim) (*model.Tokens, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}
	url := fmt.Sprintf("%s%s%s", client.Environment.BASE_URL, client.Environment.AUTH_PORT, client.Environment.GENERATE_TOKEN_URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	var tokens model.Tokens
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &tokens, nil
}

func (client *AuthClient) RefreshToken(input model.RefreshToken) (*model.Tokens, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}
	url := fmt.Sprintf("%s%s%s", client.Environment.BASE_URL, client.Environment.AUTH_PORT, client.Environment.REFRESH_TOKEN_URL)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	var tokens model.Tokens
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &tokens, nil
}

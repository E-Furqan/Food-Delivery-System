package AuthClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

func (client *AuthClient) GenerateToken(input model.UserClaim) (*model.Tokens, error) {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}
	url := fmt.Sprintf("%s%s%s", client.AuthClientEnv.BASE_URL, client.AuthClientEnv.AUTH_PORT, client.AuthClientEnv.GENERATE_TOKEN_URL)
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
	url := fmt.Sprintf("%s%s%s", client.AuthClientEnv.BASE_URL, client.AuthClientEnv.AUTH_PORT, client.AuthClientEnv.REFRESH_TOKEN_URL)
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

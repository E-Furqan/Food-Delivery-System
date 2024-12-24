package userClient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
)

func (userClient *UserClient) FetchEmail(token string) (*model.UserEmail, error) {

	jsonData, err := json.Marshal("")
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %v", err)
	}

	url, err := utils.CreateUrl(userClient.envVar.BASE_URL,
		userClient.envVar.USER_PORT,
		userClient.envVar.Fetch_email_URL)

	log.Print("URL      ", url)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	req, err := utils.CreateRequest(url, jsonData, "GET", token)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	client := utils.CreateHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update order status: received HTTP %d", resp.StatusCode)
	}

	var email model.UserEmail
	limit := int64(1 << 20)
	if err := json.NewDecoder(io.LimitReader(resp.Body, limit)).Decode(&email); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &email, nil
}

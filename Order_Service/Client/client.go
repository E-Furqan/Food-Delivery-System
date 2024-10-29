package RestaurantClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	payload "github.com/E-Furqan/Food-Delivery-System/Payload"
)

type RestaurantClient struct {
	envVar *environmentVariable.Environment
}

func NewClient(envVar *environmentVariable.Environment) *RestaurantClient {
	return &RestaurantClient{
		envVar: envVar,
	}
}

func (client *RestaurantClient) GetItems(getItems payload.GetItems) ([]payload.Items, error) {

	body, err := json.Marshal(getItems)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	url := fmt.Sprintf("%s%s%s", client.envVar.BASE_URL, client.envVar.RESTAURANT_PORT, client.envVar.Get_Items_URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
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

	var items []payload.Items
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("error un marshaling response: %v", err)
	}

	return items, nil
}

// func (client *Client) ProcessOrder(ProcessOrder payload.ProcessOrder, forUser bool) error {

// 	jsonData, err := json.Marshal(ProcessOrder)
// 	if err != nil {
// 		return fmt.Errorf("error marshaling input: %v", err)
// 	}

// 	var url string

// 	if forUser {
// 		url = fmt.Sprintf("%s%s%s", client.BaseUrl, client.USER_PORT, client.Process_Order_User_URL)
// 	} else {
// 		url = fmt.Sprintf("%s%s%s", client.BaseUrl, client.RESTAURANT_PORT, client.Process_Order_Restaurant_URL)
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

// 	if err != nil {
// 		return fmt.Errorf("error creating request: %v", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("error sending request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		body := (resp.Body)
// 		return fmt.Errorf("received non-200 response: %v, body: %s", resp.Status, body)
// 	}

// 	return nil
// }

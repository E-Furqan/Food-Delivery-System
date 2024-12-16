package RestaurantClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

func (client *RestaurantClient) GetItems(getItems model.CombineOrderItem) ([]model.Items, error) {

	body, err := json.Marshal(getItems)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	url := fmt.Sprintf("%s%s%s", client.envVar.BASE_URL, client.envVar.RESTAURANT_PORT, client.envVar.Get_Items_URL)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
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

	var items []model.Items
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("error un marshaling response: %v", err)
	}

	return items, nil
}

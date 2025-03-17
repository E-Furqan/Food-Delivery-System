package WorkFlow

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
)

func (WorkFlow *WorkFlowClient) PlaceORder(order model.CombineOrderItem, token string) error {

	jsonData, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("error marshaling input: %v", err)
	}

	url := fmt.Sprintf("%s%s%s", WorkFlow.envVar.BASE_URL, WorkFlow.envVar.WORKFLOW_PORT, WorkFlow.envVar.PLACE_ORDER_URL)

	req, err := utils.CreateAuthorizedRequest(url, jsonData, "GET", token)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	client := utils.CreateHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update order status: received HTTP %d", resp.StatusCode)
	}

	return nil
}

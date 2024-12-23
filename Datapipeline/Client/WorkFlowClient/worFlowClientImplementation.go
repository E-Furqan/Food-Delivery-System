package workflowClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

func (client *WorkflowClient) DatapipelineSync(Pipeline model.Pipeline) error {

	body, err := json.Marshal(Pipeline)
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	url := fmt.Sprintf("%s:%s%s", client.envVar.BASE_URL, client.envVar.PORT, client.envVar.DATAPIPELINE_WORKFLOW_URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	return nil
}

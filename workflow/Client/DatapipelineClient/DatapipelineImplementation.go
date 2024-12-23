package datapipelineClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

func (client *DatapipelineClient) FetchSourceConfiguration(source model.Source) (model.Config, error) {

	body, err := json.Marshal(source)
	if err != nil {
		return model.Config{}, fmt.Errorf("error marshaling request body: %v", err)
	}

	url := fmt.Sprintf("%s:%s%s", client.envVar.BASE_URL, client.envVar.DATAPIPELINE_PORT, client.envVar.FETCH_SOURCE_CONFIGURATION_URL)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
	if err != nil {
		return model.Config{}, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.Config{}, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Config{}, fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	var sourceConfig model.Config
	if err := json.NewDecoder(resp.Body).Decode(&sourceConfig); err != nil {
		return model.Config{}, fmt.Errorf("error un marshaling response: %v", err)
	}

	return sourceConfig, nil
}

func (client *DatapipelineClient) FetchDestinationConfiguration(destination model.Destination) (model.Config, error) {

	body, err := json.Marshal(destination)
	if err != nil {
		return model.Config{}, fmt.Errorf("error marshaling request body: %v", err)
	}

	url := fmt.Sprintf("%s:%s%s", client.envVar.BASE_URL, client.envVar.DATAPIPELINE_PORT, client.envVar.FETCH_DESTINATION_CONFIGURATION_URL)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
	if err != nil {
		return model.Config{}, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.Config{}, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Config{}, fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	var sourceConfig model.Config
	if err := json.NewDecoder(resp.Body).Decode(&sourceConfig); err != nil {
		return model.Config{}, fmt.Errorf("error un marshaling response: %v", err)
	}

	return sourceConfig, nil
}

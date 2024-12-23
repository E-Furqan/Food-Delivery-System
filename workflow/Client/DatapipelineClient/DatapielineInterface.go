package datapipelineClient

import model "github.com/E-Furqan/Food-Delivery-System/Models"

type DatapipelineClient struct {
	envVar model.DatapipelineClientEnv
}

func NewClient(envVar model.DatapipelineClientEnv) *DatapipelineClient {
	return &DatapipelineClient{
		envVar: envVar,
	}
}

type DatapipelineClientInterface interface {
	FetchSourceConfiguration(source model.Source) (model.Config, error)
	FetchDestinationConfiguration(destination model.Destination) (model.Config, error)
}
